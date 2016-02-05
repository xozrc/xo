package serverutil_test

import (
	"errors"
	"fmt"
	"github.com/xo/common"
	"github.com/xo/core"
	"github.com/xo/serverutil"

	"strconv"
	"sync"
	"testing"
)

var print = fmt.Print

type writeFinalHandlerTest struct {
	writeSum int64
}

func (wft *writeFinalHandlerTest) Error(c serverutil.Context, err error) {
	fmt.Printf("error:[%v]\n", err.Error())
}

func (wft *writeFinalHandlerTest) Write(c serverutil.Context, rawMsg interface{}) (msg interface{}, err error) {
	tempNum, ok := rawMsg.(int)
	if !ok {
		return nil, errors.New("wrong msg type")
	}
	wft.writeSum += int64(tempNum)

	return
}

type readFinalHandlerTest struct {
	readSum int64
}

func (rwt *readFinalHandlerTest) Error(c serverutil.Context, err error) {
	fmt.Printf("error:[%v]\n", err.Error())
}

func (rwt *readFinalHandlerTest) Read(c serverutil.Context, rawMsg interface{}) (msg interface{}, err error) {
	tempNum, ok := rawMsg.(int)
	if !ok {
		return nil, errors.New("wrong msg type")
	}
	rwt.readSum += int64(tempNum)

	return
}

func TestFinalContextResult(t *testing.T) {

	testNum := 100000
	queueSize := 1
	taskCacheSize := 1000
	rq := core.NewTaskQueue(queueSize, taskCacheSize)
	wq := core.NewTaskQueue(queueSize, taskCacheSize)

	srh := &readFinalHandlerTest{}
	swh := &writeFinalHandlerTest{}
	rc := serverutil.NewContext(nil, nil, rq, srh)
	wc := serverutil.NewContext(nil, nil, wq, swh)

	wg := sync.WaitGroup{}
	wg.Add(testNum)

	for i := 0; i < testNum; i++ {
		go func(ii int) {
			doneRead := rc.Read(ii)
			doneWrite := wc.Write(ii)
			<-doneRead
			<-doneWrite
			wg.Done()
		}(i)
	}
	wg.Wait()

	tempReadSum, tempWriteSum := sum(testNum)

	common.AssertTest(t, tempReadSum == int(srh.readSum), "final read result ["+strconv.Itoa(tempReadSum)+"] no equal context result ["+strconv.Itoa(int(srh.readSum))+"]")
	common.AssertTest(t, tempWriteSum == int(swh.writeSum), "final read result ["+strconv.Itoa(tempWriteSum)+"] no equal context result ["+strconv.Itoa(int(swh.writeSum))+"]")
}

func sum(num int) (readSum int, writeSum int) {
	for i := 0; i < num; i++ {
		readSum += i
		writeSum += i
	}
	return
}

var (
	errorRead  = errors.New("error read")
	errorWrite = errors.New("error write")
)

type errorHandlerTest struct {
	err error
}

func (eht *errorHandlerTest) Error(c serverutil.Context, err error) {
	eht.err = err
}

func (eht *errorHandlerTest) Read(c serverutil.Context, rawmsg interface{}) (msg interface{}, err error) {
	return nil, errorRead
}

func (eht *errorHandlerTest) Write(c serverutil.Context, rawmsg interface{}) (msg interface{}, err error) {
	return nil, errorWrite
}

//test error
func TestErrorContext(t *testing.T) {

	testNum := 1
	queueSize := 1
	taskCacheSize := 1
	eq := core.NewTaskQueue(queueSize, taskCacheSize)
	ewq := core.NewTaskQueue(queueSize, taskCacheSize)
	eh := &errorHandlerTest{}
	ewh := &errorHandlerTest{}

	ec := serverutil.NewContext(nil, nil, eq, eh)
	ewc := serverutil.NewContext(nil, nil, ewq, ewh)
	wg := sync.WaitGroup{}
	wg.Add(testNum)

	for i := 0; i < testNum; i++ {
		go func(ii int) {
			doneRead := ec.Read(i)
			doneWrite := ewc.Write(i)
			<-doneRead
			<-doneWrite
			wg.Done()
		}(i)
	}
	wg.Wait()

	common.AssertTest(t, eh.err == errorRead, "error  read handler")
	common.AssertTest(t, ewh.err == errorWrite, "error write handler")

}

//test panic error
type panicHandlerTest struct {
	err error
}

func (pht *panicHandlerTest) Error(c serverutil.Context, err error) {
	pht.err = err

}

func (eht *panicHandlerTest) Read(c serverutil.Context, rawmsg interface{}) (msg interface{}, err error) {
	panic("panic handler read test")
	return nil, errorRead
}

func (eht *panicHandlerTest) Write(c serverutil.Context, rawmsg interface{}) (msg interface{}, err error) {
	panic("panic handler write test")
	return nil, errorWrite
}

//test panic
func TestPanicContext(t *testing.T) {

	testNum := 1
	queueSize := 1
	taskCacheSize := 1
	prq := core.NewTaskQueue(queueSize, taskCacheSize)
	pwq := core.NewTaskQueue(queueSize, taskCacheSize)
	prh := &panicHandlerTest{}
	pwh := &panicHandlerTest{}

	prc := serverutil.NewContext(nil, nil, prq, prh)
	pwc := serverutil.NewContext(nil, nil, pwq, pwh)
	wg := sync.WaitGroup{}
	wg.Add(testNum)

	for i := 0; i < testNum; i++ {
		go func(ii int) {
			doneRead := prc.Read(i)
			doneWrite := pwc.Write(i)
			<-doneRead
			<-doneWrite
			wg.Done()
		}(i)
	}
	wg.Wait()

	common.AssertTest(t, prh.err == serverutil.ContextProcessMsgErr, "panic  read handler")
	common.AssertTest(t, pwh.err == serverutil.ContextProcessMsgErr, "panic write handler")

}

//test panic error
type headContextHandlerTest struct {
}

func (pht *headContextHandlerTest) Error(c serverutil.Context, err error) {

}

func (eht *headContextHandlerTest) Read(c serverutil.Context, rawmsg interface{}) (msg interface{}, err error) {
	msg = rawmsg
	return
}

func (eht *headContextHandlerTest) Write(c serverutil.Context, rawmsg interface{}) (msg interface{}, err error) {
	msg = rawmsg
	return
}

//some transform
