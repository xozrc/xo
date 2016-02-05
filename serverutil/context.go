package serverutil

type Context interface {
	Pipeline() Pipeline
	Prev() Context
	Next() Context
	InputHandler() InputHandler
	OutputHandler() OutputHandler
	ErrorHandler() ErrorHandler

	OnRead(rawMsg interface{}) chan bool
	OnWrite(rawMsg interface{}) chan bool
	OnError(err error) chan bool
	OnClose()
}

type context struct {
	pipeline      Pipeline
	prev          Context
	next          Context
	inputHandler  InputHandler
	outputHandler OutputHandler
	errorHandler  ErrorHandler
}

func (c *context) Pipeline() Pipeline {
	return c.pipeline
}

func (c *context) Prev() Context {
	return c.prev
}

func (c *context) Next() Context {
	return c.next
}

func (c *context) InputHandler() InputHandler {
	return c.inputHandler
}

func (c *context) OutputHandler() OutputHandler {
	return c.outputHandler
}

func (c *context) ErrorHandler() Handler {
	return c.errorHandler
}

func (c *context) OnRead(rawMsg interface{}) chan bool {
	return c.process(msg, nil, false)
}
func (c *context) OnWrite(msg interface{}) chan bool {
	return c.process(msg, nil, true)
}

func (c *context) OnWriteError(err error) chan bool {
	return c.process(nil, err, true)
}

func (c *context) OnReadError(err error) chan bool {

}

func (c *context) process(msg interface{}, err error, written bool) chan bool {
	t := &msgProcessTask{context: c, msg: msg, err: err, written: written}
	tt := core.NewTask(t)
	c.queue <- tt
	return tt.Done
}

type Pipeline interface {
	Session() Session
	Head() Context
	Tail() Context
	Write(msgbyte []byte)
	Close(err error)

	OnRead(msgByte []byte)
	OnWrite(msgByte []byte)
	OnClose(err error)
}

type pipeline struct {
	session Session
	head    Context
	tail    Context
}

func (p *pipeline) Session() Session {
	return p.session
}

func (p *pipeline) Head() Context {
	return p.head
}

func (p *pipeline) Tail() Context {
	return p.tail
}

func (p *pipeline) Close(err error) {
	p.session.Close()
	p.OnClose(err)
}

func (p *pipeline) OnRead(msgByte []byte) {
	p.head.Read(rawMsg)
}

func (p *pipeline) OnWrite(msgByte []byte) {
	p.session.Send(msgByte)
}

func (p *pipeline) OnClose(err error) {

}

type Handler interface {
	Error(c Context, err error) (err error)
}

type TrivialInputHandler interface {
	Read(c Context, rawMsg interface{}) (msg interface{}, err error)
}

type TrivialOutputHandler interface {
	Write(c Context, rawMsg interface{}) (msg interface{}, err error)
}

type InputHandler interface {
	Handler
	TrivialInputHandler
}

type OutputHandler interface {
	Handler
	TrivialOutputHandler
}

type InputOutputHandler interface {
	Handler
	TrivialInputHandler
	TrivialOutputHandler
}

type msgProcessTask struct {
	context Context
	msg     interface{}
	err     error
	written bool
}

func newMsgProcessTask(c Context, msg interface{}, err error, written bool) *msgProcessTask {
	mpt := &msgProcessTask{}
	mpt.c = c
	mpt.msg = msg
	mpt.err = err
	mpt.written = written
	return mpt
}

func (*msgProcessTask) nextContext() Context {

}

func (mpt *msgProcessTask) Run() {
	defer func() {
		if rec := recover(); rec != nil {
			logger.Printf("recover errr [%v]", rec)
			if mpt.context.Handler() == nil {
				logger.Panicf("handler is null")
				return
			}

			err := mpt.context.Handler().Error(mpt.context, ContextProcessMsgErr)
			if err != nil {
				nctx := mpt.nextContext()
				if nctx == nil {
					logger.Panicf("no catch error[%v]", err.Error())
				}
				<-nctx.OnError(err)
			}
		}
	}()

	if mpt.written {
		o, ok := mpt.context.Handler().(OutputHandler)
		if !ok {
			mpt.context.Handler().Error(mpt.context, handlerErr)
			return
		}

		msg, err := o.Write(mpt.context, mpt.msg)
		if err != nil {
			mpt.context.Handler().Error(mpt.context, err)
			return
		}

		for {
			prev := mpt.context.Prev()
			if prev == nil {
				return
			}

			if !prev.Output() {
				continue
			}

			<-prev.OnWrite(msg)
			return
		}
	} else {
		i, ok := mpt.context.Handler().(InputHandler)
		if !ok {
			mpt.context.Handler().Error(mpt.context, handlerErr)
			return
		}

		msg, err := i.Read(mpt.context, mpt.msg)
		if err != nil {
			mpt.context.Handler().Error(mpt.context, err)
			return
		}

		for {
			next := mpt.context.Next()
			if next == nil {
				return
			}

			if !next.Input() {
				continue
			}

			<-next.OnRead(msg)
			return
		}
	}
}
