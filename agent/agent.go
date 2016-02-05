package agent

import (
	"github.com/xo/cherry"
	"golang.org/x/net/context"
)

var logger = log.New(os.Stdout, "[agent]", 0)

const (
	agentServiceKey = "AgentService"
)

type AgentService struct {
}

func (as *AgentService) Init() {

}

func (as *AgentService) Start() chan error {

}

func (as *AgentService) Stop() chan error {

}

func (as *AgentService) Destroy() {

}

func RedirectHandler(ctx context.Context, msg cherry.Message) (result interface{}, err error) {
	ch := cherry.CherryValue(ctx)
	if ch == nil {
		err = cherry.CherryNilError
		return
	}
	as := ch.ServiceManager.Get(agentServiceKey)
}

type ForwardHandler struct {
}

func (fh *ForwardHandler) Handle(ctx context.Context, msgByte interface{}) (result interface{}, err error) {
	return
}
