package knet

import (
	"fmt"
	"net"
	"testing"
	tt "time"
)

type testHandler struct {
}

func (th *testHandler) OnData(c IConn) (out []byte) {
	out = c.Read()
	fmt.Println("on data", out)
	return out
}
func (th *testHandler) OnOpen(c IConn) {
	println("on open", c.Addr().String())
}
func (th *testHandler) OnClose(c IConn) {
	println("on close", c.Fd())
}

func TestKanet(t *testing.T) {
	Server(8)
	Listen(":54167", &testHandler{})

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	// <-c
	ch := make(chan bool, 100)
	for i := 0; i < cap(ch); i++ {
		go func() {
			conn, err := net.Dial("tcp", "127.0.0.1:54167")
			if err != nil {
				println("conn error!", err.Error())
			} else {
				tt.Sleep(10 * tt.Millisecond)

				n, err := conn.Write([]byte{1, 2, 3, 4, 5, 6, 67})
				if err != nil {
					println("write error", err.Error())
				} else {
					println("write success =", n)
				}
				buf := make([]byte, 1024)
				n, err = conn.Read(buf)
				if n > 0 {
					fmt.Println("read buf:", buf[:n])
				}

				tt.Sleep(100 * tt.Millisecond)
				conn.Close()
			}
			ch <- true
		}()
	}
	for i := 0; i < cap(ch); i++ {
		<-ch
	}
	Stop()
}
