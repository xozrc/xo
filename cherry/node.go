package cherry

import (
	"errors"
	"golang.org/x/net/context"
)

const (
	nodeKey        = "node"
	middleNodeKey  = "middleNode"
	prevContextKey = "prevContext"
)

const (
	queueSize = 1000
)

var (
	SessionNilError = errors.New("session null")
)

type Node interface {
	Receive(msg *Message) error
	Send(msg *Message) error
	Forward(msg *Message) error
	Reply(msg *Message) error
	Error(err error)
	Start()
	Close() error
}

func NodeValue(ctx context.Context) Node {
	n := ctx.Value(nodeKey)
	if n == nil {
		return nil
	}
	n1, _ := n.(Node)
	return n1
}

func WithNode(parent context.Context, node Node) context.Context {
	return context.WithValue(parent, nodeKey, node)
}

func PrevContextValue(ctx context.Context) context.Context {
	n := ctx.Value(prevContextKey)
	if n == nil {
		return nil
	}
	n1, _ := n.(context.Context)
	return n1
}

func WithPrevContext(parent context.Context, ctx context.Context) context.Context {
	return context.WithValue(parent, prevContextKey, ctx)
}

func NewPipelienNode(parent context.Context, rh Handler, wh Handler, fh Handler, rph Handler, next Node) Node {
	p := &PipelineNode{}
	p.ctx, _ = context.WithCancel(parent)
	p.rh = rh
	p.rq = make(chan *Message, queueSize)
	p.wh = wh
	p.wq = make(chan *Message, queueSize)
	p.fh = fh
	p.fq = make(chan *Message, queueSize)
	p.rph = rph
	p.rpq = make(chan *Message, queueSize)
	p.next = next
	return p
}

type PipelineNode struct {
	ctx  context.Context
	rq   chan *Message
	rh   Handler
	wq   chan *Message
	wh   Handler
	fq   chan *Message
	fh   Handler
	rpq  chan *Message
	rph  Handler
	next Node
}

func (dn *PipelineNode) Receive(msg *Message) error {
	dn.rq <- msg
	logger.Println("receive end")
	return nil
}

func (dn *PipelineNode) Send(msg *Message) error {

	dn.wq <- msg
	return nil
}

func (dn *PipelineNode) Reply(msg *Message) error {
	dn.rpq <- msg
	return nil
}

func (dn *PipelineNode) Forward(msg *Message) error {
	dn.fq <- msg
	return nil
}

func (dn *PipelineNode) Error(err error) {
	logger.Printf("err %s\n", err.Error())
}

func (dn *PipelineNode) Start() {
	logger.Printf("pipeline start")
	for {
		var err error
		select {
		case rm := <-dn.rq:
			err = dn.handleReceiveMsg(rm)
		case wm := <-dn.wq:
			err = dn.handleWriteMsg(wm)
		case fm := <-dn.fq:
			err = dn.handleForwardMsg(fm)
		case rpm := <-dn.rpq:
			err = dn.handleReplyMsg(rpm)
		}

		if err != nil {
			dn.Error(err)
		}

	}
	logger.Printf("pipeline end")
}

func (dn *PipelineNode) handleReceiveMsg(msg *Message) (err error) {
	logger.Printf("rece")
	tempCtx := WithPrevContext(context.Background(), msg.Ctx)
	tempCtx = WithNode(tempCtx, dn)

	result, err := dn.rh.Handle(tempCtx, msg.Msg)
	if err != nil {
		return
	}
	logger.Println(result)
	if result != nil {
		tc := &Message{}
		tc.Ctx = tempCtx
		tc.Msg = result
		err = dn.Forward(tc)
	}
	return
}

func (dn *PipelineNode) handleWriteMsg(wm *Message) (err error) {
	logger.Println("wr")
	var result interface{}
	if dn.wh == nil {
		result = wm.Msg
	} else {
		result, err = dn.wh.Handle(wm.Ctx, wm.Msg)
		if err != nil {
			return
		}
	}

	if result != nil {
		pctx := PrevContextValue(dn.ctx)
		if pctx == nil {
			sess := SessionValue(dn.ctx)
			if sess == nil {
				err = SessionNilError
				return
			}
			content, ok := wm.Msg.([]byte)
			if !ok {
				logger.Println("write error")
			}

			err = sess.Send(content)
			logger.Println("wr done")
			return err
		}
		pNode := NodeValue(pctx)
		tm := &Message{}
		tm.Ctx = dn.ctx
		tm.Msg = result
		err = pNode.Reply(tm)
	}
	logger.Printf("wr done")
	return
}

func (dn *PipelineNode) handleForwardMsg(fm *Message) (err error) {
	logger.Printf("forward")
	result, err := dn.fh.Handle(fm.Ctx, fm.Msg)
	if err != nil {
		return
	}

	if result != nil {
		tm := &Message{}
		tm.Ctx = dn.ctx
		tm.Msg = result
		err = dn.next.Receive(tm)
	}
	return
}

func (dn *PipelineNode) handleReplyMsg(rpm *Message) (err error) {
	logger.Printf("reply")
	cctx := PrevContextValue(dn.ctx)

	result, err := dn.rph.Handle(cctx, rpm.Msg)
	if err != nil {
		return
	}
	if result != nil {
		n := NodeValue(cctx)
		tm := &Message{}
		tm.Ctx = cctx
		tm.Msg = result
		n.Send(tm)
	}
	return
}

func (dn *PipelineNode) Close() error {
	close(dn.rq)
	close(dn.wq)
	close(dn.fq)
	close(dn.rpq)
	return nil
}
