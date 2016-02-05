package core

type Initer interface {
	Init()
}

type Destroyer interface {
	Destroy()
}

type Lifecycle interface {
	Initer
	Destroyer
}
