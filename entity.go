package xo

import (
	"errors"
	"reflect"
)

type Entity interface {
	Init()
	Destroy()
	GetComponent(ty reflect.Type) Component
	AddComponent(com Component)
}

type entity struct {
	components map[reflect.Type]Component
}

//init
func (e *entity) Init() {
	for _, com := range e.components {
		com.Init()
	}
}

//destroy
func (e *entity) Destroy() {
	for _, com := range e.components {
		com.Destroy()
	}
}

func (e *entity) GetComponent(ty reflect.Type) (com Component) {
	com = e.components[ty]
	return
}

func (e *entity) AddComponent(com Component) (err error) {
	ty := reflect.TypeOf(com)
	com2 := e.components[ty]
	if com2 != nil {
		err = errors.New("component alreay exist")
		return
	}
	e.components[ty] = com
	return
}

func NewEntity() Entity {
	return &entity{components: make(map[reflect.Type]Component)}
}
