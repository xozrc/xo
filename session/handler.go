package session

import (
	"golang.org/x/net/context"
)

type HandlerFunc func(ctx context.Context, session Session, msgByte interface{}) (error, interface{})

func (hf HandlerFunc) Handle(ctx context.Context, session Session, msgByte interface{}) (error, interface{}) {
	return hf(ctx, session, msgByte)
}

type Handler interface {
	Handle(ctx context.Context, session Session, msgByte interface{}) (error, interface{})
}
