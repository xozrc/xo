package serverutil

import (
	"errors"
	"fmt"
	"github.com/xo/core"
)

const (
	defaultTaskCacheSize = 1000
)

var (
	handlerErr           = errors.New("handler error")
	ContextProcessMsgErr = errors.New("process msg error")
	defauleQueue         = core.TaskQueue(make(chan *core.Task, defaultTaskCacheSize))
)

type Handler interface {
	Error(c Context, err error)
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
	written bool
}

func newMsgProcessTask(c Context, msg interface{}, written bool) *msgProcessTask {
	return &msgProcessTask{context: c, msg: msg, written: written}
}

func (mpt *msgProcessTask) Run() {
	defer func() {
		if rec := recover(); rec != nil {
			logger.Printf("recover errr [%v]", rec)
			if mpt.context.Handler() == nil {
				logger.Panicf("handler is null")
				return
			}

			mpt.context.Handler().Error(mpt.context, ContextProcessMsgErr)
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

			prev.Write(msg)
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

			next.Read(msg)
			return
		}
	}
}

type Context interface {
	Prev() Context
	Next() Context
	SetNext(c Context)
	Handler() Handler
	Write(msg interface{}) chan bool
	Read(msg interface{}) chan bool
	Input() bool
	Output() bool
}

func NewContext(prev Context, next Context, queue core.TaskQueue, h Handler) Context {
	ctx := &context{}
	ctx.prev = prev
	ctx.next = next
	ctx.queue = queue
	ctx.h = h
	return ctx
}

type context struct {
	prev  Context
	next  Context
	queue core.TaskQueue
	h     Handler
}

func (c *context) Prev() Context {
	return c.prev
}

func (c *context) Next() Context {
	return c.next
}

func (c *context) SetNext(ctx Context) {
	c.next = ctx
}

func (c *context) Handler() Handler {
	return c.h
}

func (c *context) Input() bool {
	_, ok := c.Handler().(InputHandler)
	return ok
}

func (c *context) Output() bool {
	_, ok := c.Handler().(OutputHandler)
	return ok
}

func (c *context) Read(msg interface{}) chan bool {
	return c.process(msg, false)
}

func (c *context) Write(msg interface{}) chan bool {
	return c.process(msg, true)
}

func (c *context) process(msg interface{}, written bool) chan bool {
	t := &msgProcessTask{context: c, msg: msg, written: written}
	tt := core.NewTask(t)
	c.queue <- tt
	return tt.Done
}

type headHandler struct {
}

func (hh *headHandler) Error(c Context, err error) {

}

func (hh *headHandler) Read(c Context, rawMsg interface{}) (msg interface{}, err error) {

}

func (hh *headHandler) Write(c Context, rawMsg interface{}) (msg interface{}, err error) {

}

type HeadContext struct {
	Context
	session Session
}

func (hc *HeadContext) Write(msg interface{}) chan bool {
	content, _ := msg.([]byte)
	hc.session.Send(content)
	done := make(chan bool, 1)
	done <- true
	return done
}

func (hc *HeadContext) AddHandler(h Handler, queue core.TaskQueue) {
	for tc := hc.Context; tc != nil; tc = tc.Next() {
		if tc.Next() == nil {
			cctx := NewContext(tc, nil, queue, h)
			tc.SetNext(cctx)
			break
		}
	}
}

func (hc *HeadContext) RemoveHandler(h Handler) {
	for tc := hc.Context; tc != nil; tc = tc.Next() {
		if tc.Handler() == h {
			tc.Prev().SetNext(nil)
			break
		}
	}

}

func NewHeadContext(queue core.TaskQueue, h Handler, session Session) *HeadContext {
	hc := &HeadContext{}
	c := NewContext(nil, nil, queue, h)
	hc.Context = c
	hc.session = session
	return hc
}

type Session interface {
	GetId() int64
	GetContext() *HeadContext

	Send(msgByte []byte)
	Close()

	OnReceive(msgByte []byte)
	OnError(err error)
	OnClose()
}

type session struct {
	id     int64
	ctx    *HeadContext
	closed bool
}

func (s *session) GetId() int64 {
	return s.id
}

func (s *session) GetContext() *HeadContext {
	return s.ctx
}

func (s *session) Send(msgByte []byte) {
	if s.closed {
		return
	}
	fmt.Printf("send msg [%v]\n", msgByte)
}

func (s *session) Close() {
	if s.closed {
		return
	}
	s.closed = true
	fmt.Printf("close session")
}

func (s *session) OnClose() {

}

func (s *session) OnReceive(msgByte []byte) {
	s.ctx.Read(msgByte)
}

func (s *session) OnError(err error) {

}

func NewSession(id int64, ctx *HeadContext) Session {
	s := &session{}
	s.ctx = ctx
	s.id = id
	return s
}

// func NewSession(c Context) Session {
// 	s := &session{}
// 	s.Context = c
// 	s.id = 1
// 	return s
// }

// type WSSession struct {
// 	Id   int64
// 	Conn *websocket.Conn
// 	*Context
// }

// func (wss *WSSession) GetId() int64 {
// 	return wss.ID
// }

// func (wss *WSSession) OnReceive(msgByte []byte) {
// 	wss.ctx.Read(msgByte)
// }

// func (wss *wsSession) Send(msgByte []byte) {
// 	_, err := wss.Conn.Write(msgByte)
// 	if err != nil {

// 	}
// }

// func (wss *wsSession) Close() {
// 	err := wss.Conn.Close()
// 	if err != nil {
// 		logger.Printf("session [%v] close with err [%v]", wss.GetId(), err.Error())
// 		return
// 	}
// 	logger.Printf("session [%v] active close", wss.GetId())
// }

// func (wss *wsSession) AddHandler(h Handler, queue core.TaskQueue) {

// }

// func (wss *wsSession) RemoveHandler(h Handler) {

// }

// func newWsSession(id int64, conn *websocket.Conn) (wss *wsSession) {
// 	wss = &wsSession{}
// 	wss.id = id
// 	wss.Conn = conn
// 	wss.ctx = NewHeadContext(defauleQueue, &DefaultHandler{}, wss)
// 	return
// }

// type DefaultHandler struct {
// }

// func (dh *DefaultHandler) Error(c Context, err error) {
// 	panic(err.Error())
// }

// func (dh *DefaultHandler) Read(c Context, rawMsg interface{}) (msg interface{}, err error) {
// 	return
// }

// func (dh *DefaultHandler) Write(c Context, rawMsg interface{}) (msg interface{}, err error) {
// 	return
// }
