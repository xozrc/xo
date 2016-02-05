package cherry

import (
	"golang.org/x/net/context"
)

type Message struct {
	Ctx context.Context
	Msg interface{}
}

func NewMessage(ctx context.Context, msg interface{}) *Message {
	m := &Message{}
	m.Ctx = ctx
	m.Msg = msg
	return m
}

type HandlerFunc func(ctx context.Context, msgByte interface{}) (interface{}, error)

func (hf HandlerFunc) Handle(ctx context.Context, msgByte interface{}) (interface{}, error) {
	return hf(ctx, msgByte)
}

type Handler interface {
	Handle(ctx context.Context, msgByte interface{}) (interface{}, error)
}

func NewPipelineHandler(hs ...Handler) Handler {
	p := &PipelineHandler{}
	for _, h := range hs {
		p.handlers = append(p.handlers, h)
	}
	return p
}

type PipelineHandler struct {
	handlers []Handler
}

func (p *PipelineHandler) Handle(ctx context.Context, msgByte interface{}) (result interface{}, err error) {
	result = msgByte
	for _, h := range p.handlers {
		result, err = h.Handle(ctx, result)
		if err != nil {
			return
		}
	}
	return
}
