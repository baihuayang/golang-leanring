package knet

import (
	"knet/netpoll"
	"knet/ringbuffer"

	"golang.org/x/sys/unix"
)

func (svr *server) acceptNewConnection(fd int) error {
	nfd, sa, err := unix.Accept(fd)
	if err != nil {
		if err == unix.EAGAIN {
			return nil
		}
		return err
	}
	if err := unix.SetNonblock(nfd, true); err != nil {
		return err
	}
	// log.Println("accpet new fd:", nfd)
	l := svr.subLoopGroup.next()
	c := &connection{
		fd:         nfd,
		sa:         sa,
		lp:         l,
		lfd:        fd,
		outBuffer:  ringbuffer.New(1024),
		inBuffer:   ringbuffer.New(1024),
		remoteAddr: netpoll.SockaddrToTCPOrUnixAddr(sa),
	}
	if ln, ok := svr.listeners[fd]; ok {
		ln.handler.OnOpen(c)
	}
	l.poller.Trigger(func() error {
		if err := l.poller.AddRead(nfd); err != nil {
			return err
		}
		l.conns[nfd] = c
		return nil
	})
	return nil
}
