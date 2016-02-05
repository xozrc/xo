package message

type MessageService struct {
	msgManager *messageManager
}

func (ms *MessageService) Init() {

}

func (ms *MessageService) Tick() {
	for {
		msg := ms.msgManager.Consume()
		msg.Handle()
	}
}

func (ms *MessageService) Destroy() {

}
