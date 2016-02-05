package core_test

import (
	"fmt"
	"github.com/xo/common"
	"github.com/xo/core"
	"sync"
	"testing"
	"time"
)

func funcTest() {
	fmt.Println("func test succ")
}

func TestSingleQueueSingleTaskCache(t *testing.T) {
	taskQueue := core.NewTaskQueue(1, 1)
	for i := 0; i < 1; i++ {
		taskQueue <- core.NewFuncTask(funcTest)
	}
}

type oneArgumentSyncOperationTest struct {
	sum     *int
	counter int
	t       *testing.T
}

func (o *oneArgumentSyncOperationTest) Run() {
	time.Sleep(10)
	expected := 0
	for i := 0; i <= o.counter; i++ {
		expected += i
	}
	*o.sum += o.counter
	fmt.Printf("input %v,get %v\n", o.counter, *o.sum)
	common.AssertTest(o.t, expected == *o.sum, fmt.Sprintf("input %v,expected %v,get %v", o.counter, expected, *o.sum))
}

func newOneArgumentSyncTask(couter int, sum *int, t *testing.T) *core.Task {
	o := &oneArgumentSyncOperationTest{counter: couter, sum: sum, t: t}
	return core.NewTask(o)
}

func TestSingleQueueMultiTaskCacheWithArg(t *testing.T) {
	sum := 0
	testNum := 1000
	taskQueue := core.NewTaskQueue(1, testNum)
	tempSync := &sync.WaitGroup{}
	tempSync.Add(testNum)

	go func() {
		for i := 0; i < testNum; i++ {
			tempTask := newOneArgumentSyncTask(i, &sum, t)
			taskQueue <- tempTask
			<-tempTask.Done
			tempSync.Done()
		}
	}()

	tempSync.Wait()
}

type argumentAsyncOperationTest struct {
	result  chan int
	counter int
}

func (o *argumentAsyncOperationTest) Run() {
	o.result <- o.counter * 2
}

func newArgumentAsyncTask(couter int, result chan int) *core.Task {
	o := &argumentAsyncOperationTest{counter: couter, result: result}
	return core.NewTask(o)
}

func TestMultQueueMultiTaskCacheWithArg(t *testing.T) {
	testQueue := 10
	tempResult := make(chan int, testQueue)

	testNum := 100000
	taskQueue := core.NewTaskQueue(testQueue, testNum)

	tempSync := &sync.WaitGroup{}
	tempSync.Add(testNum)
	go func() {
		for i := 0; i < testNum; i++ {
			tempO := newArgumentAsyncTask(i, tempResult)
			taskQueue <- tempO
			<-tempO.Done
			tempSync.Done()
		}
	}()

	cursor := 0
	sum := 0
	for i := range tempResult {
		sum += i
		cursor += 1
		if cursor == testNum {
			break
		}
	}

	expeceted := 0
	for i := 0; i < testNum; i++ {
		expeceted += i
	}
	expeceted *= 2
	common.AssertTest(t, expeceted == sum, fmt.Sprintf("expected %v,get %v", expeceted, sum))
}

func BenchmarkSingleTaskQueue(b *testing.B) {
	taskQueue := core.NewTaskQueue(1, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		taskQueue <- core.NewFuncTask(funcTest)
	}
}
