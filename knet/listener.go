package knet

import (
	"net"
	"os"

	"golang.org/x/sys/unix"
)

type listener struct {
	ln      net.Listener
	lnAddr  net.Addr
	network string
	fd      int
	f       *os.File
	handler eventHandler
	// react func(c connection) error
}

type eventHandler interface {
	OnData(c IConn) (out []byte)
	OnOpen(c IConn)
	OnClose(c IConn)
}

func (ln *listener) close() {
	if ln.f != nil {
		printIfError(ln.f.Close())
	}
	if ln.ln != nil {
		printIfError(ln.ln.Close())
	}
}

func (ln *listener) system() error {
	var err error
	switch netln := ln.ln.(type) {
	case nil:
	case *net.TCPListener:
		ln.f, err = netln.File()
	case *net.UnixListener:
		ln.f, err = netln.File()
	}
	if err != nil {
		ln.close()
		return err
	}
	ln.fd = int(ln.f.Fd())
	ln.lnAddr = ln.ln.Addr()
	return unix.SetNonblock(ln.fd, true)
}
