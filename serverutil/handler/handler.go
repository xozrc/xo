package handler

type ReadHandler interface {
	Read(c Context, obj interface{})
}

type WriteHandler interface {
	Write()
}

type readHandler struct {
}
