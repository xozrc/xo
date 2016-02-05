package cherry_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xo/cherry"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
	"net"
	"net/http"

	"strings"
	"sync"

	"testing"
)

var once sync.Once
var serverAddr string
var wsc *cherry.Cherry

func startCherryServer() {
	wsc = cherry.NewWSCherry()
	wsc.Init()
	serverAddr = wsc.RunTest()
}

func newConfig(t *testing.T, path string) *websocket.Config {
	config, _ := websocket.NewConfig(fmt.Sprintf("ws://%s%s", serverAddr, path), "http://localhost")
	return config
}

func echoHandler(ctx context.Context, msg interface{}) (result interface{}, err error) {
	n := cherry.NodeValue(ctx)
	n.Send(cherry.NewMessage(ctx, msg))
	return
}

func echoServer(conn *websocket.Conn) {
	sess := cherry.NewWSSession(wsc.Context(), conn, cherry.HandlerFunc(echoHandler), nil, nil, nil, nil)
	cherry.WSServe(sess, conn)
}

type Count struct {
	S string
	N int
}

type CounterHandler struct {
	C *Count
}

func (ch *CounterHandler) Handle(ctx context.Context, msg interface{}) (result interface{}, err error) {

	err = json.Unmarshal(msg.([]byte), &ch.C)
	if err != nil {
		return
	}
	ch.C.N++
	ch.C.S = strings.Repeat(ch.C.S, ch.C.N)
	contentByte, err := json.Marshal(ch.C)

	if err != nil {
		return
	}
	n := cherry.NodeValue(ctx)
	n.Send(cherry.NewMessage(ctx, contentByte))
	return
}

func countServer(conn *websocket.Conn) {
	sess := cherry.NewWSSession(wsc.Context(), conn, &CounterHandler{}, nil, nil, nil, nil)
	cherry.WSServe(sess, conn)
}

func TestEcho(t *testing.T) {
	http.Handle("/echo", websocket.Handler(echoServer))
	once.Do(startCherryServer)

	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}

	conn, err := websocket.NewClient(newConfig(t, "/echo"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}

	msg := []byte("hello, world\n")
	if _, err := conn.Write(msg); err != nil {
		t.Errorf("Write: %v", err)
	}
	var actual_msg = make([]byte, 512)
	n, err := conn.Read(actual_msg)
	if err != nil {
		t.Errorf("Read: %v", err)
	}
	actual_msg = actual_msg[0:n]
	if !bytes.Equal(msg, actual_msg) {
		t.Errorf("Echo: expected %q got %q", msg, actual_msg)
	}
	conn.Close()
}

func TestCount(t *testing.T) {
	http.Handle("/count", websocket.Handler(countServer))
	once.Do(startCherryServer)

	// websocket.Dial()
	client, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	conn, err := websocket.NewClient(newConfig(t, "/count"), client)
	if err != nil {
		t.Errorf("WebSocket handshake error: %v", err)
		return
	}

	var count Count
	count.S = "hello"
	if err := websocket.JSON.Send(conn, count); err != nil {
		t.Errorf("Write: %v", err)
	}

	if err := websocket.JSON.Receive(conn, &count); err != nil {
		t.Errorf("Read: %v", err)
	}
	if count.N != 1 {
		t.Errorf("count: expected %d got %d", 1, count.N)
	}
	if count.S != "hello" {
		t.Errorf("count: expected %q got %q", "hello", count.S)
	}
	if err := websocket.JSON.Send(conn, count); err != nil {
		t.Errorf("Write: %v", err)
	}
	if err := websocket.JSON.Receive(conn, &count); err != nil {
		t.Errorf("Read: %v", err)
	}
	if count.N != 2 {
		t.Errorf("count: expected %d got %d", 2, count.N)
	}
	if count.S != "hellohello" {
		t.Errorf("count: expected %q got %q", "hellohello", count.S)
	}

	conn.Close()
}
