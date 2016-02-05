package cherry

import (
	"github.com/xo/session"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

const (
	sessionKey = "session"
)

func SessionValue(ctx context.Context) session.Session {
	ts := ctx.Value(sessionKey)
	if ts == nil {
		return nil
	}
	ts1, _ := ts.(session.Session)
	return ts1
}

func WithSession(parent context.Context, sess session.Session) context.Context {
	return context.WithValue(parent, sessionKey, sess)
}

func NewWSSession(ctx context.Context, conn *websocket.Conn, rh Handler, wh Handler, fh Handler, rph Handler, next Node) session.Session {
	wss := &WSSession{}
	sessCtx := WithSession(ctx, wss)
	sn := NewPipelienNode(sessCtx, rh, wh, fh, rph, next)
	go func(n Node) {
		n.Start()
	}(sn)
	wss.ctx, wss.cancelFunc = context.WithCancel(ctx)
	wss.conn = conn
	wss.Node = sn
	return wss
}

type WSSession struct {
	Node
	conn       *websocket.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (s *WSSession) GetId() int64 {
	return 1
}

func (s *WSSession) OnOpen() {

}

func (s *WSSession) OnReceive(msgByte []byte) {
	mess := NewMessage(s.ctx, msgByte)
	err := s.Node.Receive(mess)
	if err != nil {
		s.Close()
		s.OnError(err)
	}
}

func (s *WSSession) OnClose() {
	s.cancelFunc()
}

func (s *WSSession) OnError(err error) {

}

func (s *WSSession) Send(msgByte []byte) error {
	return websocket.Message.Send(s.conn, msgByte)
}

func (s *WSSession) Close() error {
	return s.conn.Close()
}
