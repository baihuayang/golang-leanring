package knet

import (
	"net"
	"testing"
)

func TestConn(t *testing.T) {
	// tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:54167")
	// if err != nil {
	// 	return
	// }

	conn, err := net.Dial("tcp", "127.0.0.1:54167")
	if err != nil {
		println("conn error!", err.Error())
		return
	}
	println(conn.LocalAddr().String(), conn.RemoteAddr().String())
	conn.Close()
	println("=======finish========")
}
