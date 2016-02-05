package cherry

import (
	"errors"
	"github.com/xo/core"
	"golang.org/x/net/context"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[cherry]", 0)

const (
	cherryKey string = "Cherry"
)

var (
	CherryNilError = errors.New("cherry is nil")
)

func CherryValue(ctx context.Context) *Cherry {
	c := ctx.Value(cherryKey)
	if c == nil {
		return nil
	}
	c1, _ := c.(*Cherry)
	return c1
}

func WithCherry(parent context.Context, c *Cherry) context.Context {
	return context.WithValue(parent, cherryKey, c)
}

type CherryServer interface {
	Context() context.Context
	Init()
	Run() error
	RunOnAddr(addr string) error
	RunTest() string
}

func NewCherry() *Cherry {
	c := &Cherry{}
	c.ServiceManager = core.NewServiceManager()
	bc := context.Background()
	c.ctx = WithCherry(bc, c)
	return c
}

//cherry deamon
type Cherry struct {
	*core.ServiceManager
	CherryServer
	ctx context.Context
}

func (c *Cherry) Init() {
	for _, v := range c.ServiceManager.Services() {
		v.Init()
	}
}

func (c *Cherry) Start() {
	for _, v := range c.ServiceManager.Services() {
		err := <-v.Start()
		if err != nil {
			c.error(err)
		}
	}
}

func (c *Cherry) Stop() {
	for _, v := range c.ServiceManager.Services() {
		err := <-v.Stop()
		if err != nil {
			c.error(err)
		}
	}
}

func (c *Cherry) Destroy() {
	for _, v := range c.ServiceManager.Services() {
		v.Destroy()
	}
}

func (c *Cherry) Context() context.Context {
	return c.ctx
}

func (c *Cherry) error(err error) {
	logger.Fatalln(err.Error())
}
