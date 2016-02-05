package message

import (
	"github.com/golang/protobuf/proto"
)

type Message interface {
	Handle()
}

type MessageManager interface {
	Produce(msg Message)
	Consume() (msg Message)
}

type messageManager struct {
	messageCh chan Message
}

func (m *messageManager) Produce(msg Message) {
	m.messageCh <- msg

}

func (m *messageManager) Consume() (msg Message) {
	return <-m.messageCh
}
