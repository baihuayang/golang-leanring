package knet

import (
	"knet/netpoll"
	"knet/ringbuffer"

	"golang.org/x/sys/unix"
)

type loop struct {
	idx    int
	poller *netpoll.Poller
	packet []byte
	conns  map[int]*connection
	svr    *server
}

func (l *loop) loopRun() {
	defer l.svr.signalShutdown()
	l.poller.Polling(l.eventHandler)
}

func (l *loop) loopAccept(fd int) error {
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

	c := &connection{
		fd:         nfd,
		lfd:        fd,
		sa:         sa,
		remoteAddr: netpoll.SockaddrToTCPOrUnixAddr(sa),
		// inBuffer:   ringbuffer.New(1024),
		outBuffer: ringbuffer.New(1024),
	}

	if err = l.poller.AddRead(c.fd); err != nil {
		return err
	}

	l.conns[c.fd] = c

	return nil
}

func (l *loop) loopIn(c *connection) error {
	n, err := unix.Read(c.fd, l.packet)
	if n == 0 || err != nil {
		if err == unix.EAGAIN {
			return nil
		}
		return l.loopCloseConn(c)
	}
	// _, _ = c.inBuffer.Write(l.packet[:n])
	c.extra = l.packet[:n]
	ln, ok := l.svr.listeners[c.lfd]
	if ok {
		out := ln.handler.OnData(c)
		if out != nil {
			c.Write(out)
		}
	}
	c.inBuffer.Write(c.extra)
	return nil
}

func (l *loop) loopOut(c *connection) error {
	top, tail := c.outBuffer.PreReadAll()
	n, err := unix.Write(c.fd, top)
	if err != nil {
		if err == unix.EAGAIN {
			return nil
		}
		return l.loopCloseConn(c)
	}
	c.outBuffer.Shift(n)
	if len(top) == n && tail != nil {
		n, err = unix.Write(c.fd, tail)
		if err != nil {
			if err == unix.EAGAIN {
				return nil
			}
			return l.loopCloseConn(c)
		}
		c.outBuffer.Shift(n)
	}
	if c.outBuffer.IsEmpty() {
		l.poller.ModRead(c.fd)
	}

	return nil
}

func (l *loop) loopCloseConn(c *connection) error {
	// lfd := c.lfd
	if err := l.poller.Delete(c.fd); err == nil {
		delete(l.conns, c.fd)
	}
	if ln, ok := l.svr.listeners[c.lfd]; ok {
		ln.handler.OnClose(c)
	}
	return nil
}
