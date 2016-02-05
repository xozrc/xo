package cherry

import (
	"github.com/xo/session"
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
	"net/http/httptest"
)

const (
	port = ":4000"
)

func NewWSCherryServer(parentCtx context.Context) CherryServer {
	cs := &cherryWSServer{}
	cs.ctx, cs.cancelFunc = context.WithCancel(parentCtx)
	return cs
}

type cherryWSServer struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (cs *cherryWSServer) Context() context.Context {
	return cs.ctx
}

func (cs *cherryWSServer) Init() {

}

func (cs *cherryWSServer) Run() error {
	return cs.RunOnAddr(port)
}

func (cs *cherryWSServer) RunOnAddr(addr string) error {
	logger.Printf("try to listen %s\n", addr)

	srv := &http.Server{Addr: addr}
	srv.
	err := srv.ListenAndServe()
	return err
}

//for test
func (cs *cherryWSServer) RunTest() string {

	server := httptest.NewServer(nil)
	serverAddr := server.Listener.Addr().String()
	logger.Print("Test Cherry WebSocket server listening on ", serverAddr)

	return serverAddr
}

func NewWSCherry() *Cherry {
	c := NewCherry()
	c.CherryServer = NewWSCherryServer(c.Context())
	c.Init()
	return c
}

func WSServe(sess session.Session, conn *websocket.Conn) {
	logger.Printf("session ip[%s] start serve\n", conn.RemoteAddr().String())
	sess.OnOpen()
	for {
		var tempContent []byte
		if err := websocket.Message.Receive(conn, &tempContent); err != nil {

			if err != io.EOF {
				logger.Printf("session ip[%s] handle data error[%s]", conn.RemoteAddr().String(), err.Error())
				sess.OnError(err)
				err = sess.Close()
				if err != nil {
					logger.Printf("session ip[%s] close error[%s]", conn.RemoteAddr().String(), err.Error())
				}
			} else {
				logger.Printf("remote ip [%s] active close\n", conn.RemoteAddr().String())
				sess.OnClose()
			}
			break
		}
		logger.Printf("receive data from %s,data [%v]\n", conn.RemoteAddr().String(), tempContent)
		sess.OnReceive(tempContent)
	}
}
