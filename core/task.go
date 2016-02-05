package core

import (
	"strconv"
)

//operation
type Operation interface {
	Run()
}

//simplest operation
type operation func()

func (o operation) Run() {
	o()
}

//task
type Task struct {
	Done      chan bool
	Operation Operation
}

func NewFuncTask(f func()) *Task {
	return NewTask(operation(f))
}

//new simplest task
func NewTask(o Operation) *Task {
	t := &Task{}
	t.Done = make(chan bool, 1)
	t.Operation = o
	return t
}

//task queue is chan
type TaskQueue chan *Task

//interface method
func (tq TaskQueue) run(size int) {
	for i := 0; i < size; i++ {
		go func() {
			for {
				t := <-tq
				t.Operation.Run()
				t.Done <- true
			}
		}()
	}
}

func NewTaskQueue(size int, cacheSize int) TaskQueue {
	if size <= int(0) {
		panic("task queue size should more than zero " + strconv.Itoa(size))
	}

	if cacheSize < int(0) {
		panic("cache task size should no less than zero " + strconv.Itoa(size))
	}

	tempOp := TaskQueue(make(chan *Task, cacheSize))
	tempOp.run(size)
	return tempOp
}
