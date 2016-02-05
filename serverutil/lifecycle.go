package serverutil

type Initer interface {
	Init()
}

type Destroyer interface {
	Destroy()
}

type Tickabler interface {
	Tick()
}

type Lifecycle interface {
	Initer
	Destroyer
	Tickabler
}
