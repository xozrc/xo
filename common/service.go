package common

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[Common]", 0)

type Service interface {
	Init()
	AfterInit()
	BeforeDestroy()
	Destroy()
}
