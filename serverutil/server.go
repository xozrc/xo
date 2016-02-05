package serverutil

import (
	"fmt"
	"github.com/go-martini/martini"

	"golang.org/x/net/websocket"

	"log"

	"os"
	"sync"
)

var logger = log.New(os.Stdout, "[Server]", 0)

func newSessionManager() (sm *SessionManager) {
	sm = &SessionManager{}
	sm.sessions = make(map[int64]Session)
	sm.rwm = &sync.RWMutex{}
	return
}

type SessionManager struct {
	rwm      *sync.RWMutex
	sessions map[int64]Session
}

func (sm *SessionManager) PutSession(s Session) bool {
	sm.rwm.Lock()
	defer sm.rwm.Unlock()
	_, ok := sm.sessions[s.GetId()]
	if ok {
		return false
	}

	sm.sessions[s.GetId()] = s
	return true
}

func (sm *SessionManager) RemoveSession(s Session) {
	sm.rwm.Lock()
	defer sm.rwm.Unlock()
	_, ok := sm.sessions[s.GetId()]
	if !ok {
		return
	}
	delete(sm.sessions, s.GetId())
}

func (sm *SessionManager) SessionById(id int64) Session {
	sm.rwm.Lock()
	defer sm.rwm.Unlock()
	s, _ := sm.sessions[id]
	return s
}

type ServerConfig struct {
	Host string
	Port int32
}

func NewWSServer(serverConfig *ServerConfig) Server {
	tempS := &wsServer{
		SessionManager: newSessionManager(),
		ServerConfig:   serverConfig,
	}
	return tempS
}

type Server interface {
	Lifecycle
	Start()
	Stop()
}

type wsServer struct {
	*SessionManager
	ServerConfig *ServerConfig

	martini *martini.ClassicMartini
}

func (wss *wsServer) Init() {
	wss.martini = martini.Classic()
	wss.martini.Any("/msg", websocket.Handler(wss.connectionOpen).ServeHTTP)
}

func (wss *wsServer) Tick() {

}

func (wss *wsServer) Destroy() {

}

func (wss *wsServer) Start() {

	addr := fmt.Sprintf("%s:%d", wss.ServerConfig.Host, wss.ServerConfig.Port)
	logger.Println("server start on " + addr)
	wss.martini.RunOnAddr(addr)

}

func (wss *wsServer) Stop() {

}

//call back
func (wss *wsServer) connectionOpen(ws *websocket.Conn) {

	// tempId := int64(0)
	// tempSession := newWsSession(tempId, ws)
	// //cache session

	// flag := wss.PutSession(tempSession)
	// if !flag {
	// 	return
	// }

	// wss.sessionOpen(tempSession)

	// for {
	// 	var tempContent []byte
	// 	if err := websocket.Message.Receive(ws, &tempContent); err != nil {
	// 		if err != io.EOF {
	// 			wss.sessionError(tempSession, err)
	// 		} else {
	// 			//remove session
	// 			wss.sessionClose(tempSession)
	// 		}
	// 		wss.RemoveSession(tempSession)
	// 		break
	// 	}

	// 	wss.sessionReceive(tempSession, tempContent)

	// }

}

func (wss *wsServer) sessionOpen(s Session) {
	logger.Printf("open session[%v]", s.GetId())
}

func (wss *wsServer) sessionReceive(s Session, content []byte) {
	//put msg in chan
	logger.Println("receive:" + string(content))
}

func (wss *wsServer) sessionClose(s Session) {
	logger.Printf("close session[%v]", s.GetId())
}

func (wss *wsServer) sessionError(s Session, err error) {
	logger.Printf("session err[%v]", err.Error())
}
