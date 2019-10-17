package knet

import (
	"knet/ringbuffer"
	"net"

	"golang.org/x/sys/unix"
)

type connection struct {
	fd         int
	lfd        int
	sa         unix.Sockaddr
	opend      bool
	lp         *loop
	localAddr  net.Addr
	remoteAddr net.Addr
	extra      []byte
	inBuffer   *ringbuffer.RingBuffer
	outBuffer  *ringbuffer.RingBuffer
}

type IConn interface {
	Read() []byte
	Addr() net.Addr
	Write(buf []byte)
	Fd() int
}

func (c *connection) Read() []byte {
	// return c.inBuffer.WithBytes(c.extra)
	if c.inBuffer.IsEmpty() {
		return c.extra
	}
	top, _ := c.inBuffer.PreReadAll()
	return append(top, c.extra...)
}

func (c *connection) Addr() net.Addr {
	return c.remoteAddr
}

func (c *connection) Write(buf []byte) {
	if !c.outBuffer.IsEmpty() {
		_, _ = c.outBuffer.Write(buf)
		return
	}
	n, err := unix.Write(c.fd, buf)
	if err != nil {
		if err == unix.EAGAIN {
			_, _ = c.outBuffer.Write(buf)
			c.lp.poller.ModReadWrite(c.fd)
			return
		}
		_ = c.lp.loopCloseConn(c)
		return
	}
	if n < len(buf) {
		_, _ = c.outBuffer.Write(buf[n:])
		c.lp.poller.ModReadWrite(c.fd)
	}
}

func (c *connection) Fd() int {
	return c.fd
}
