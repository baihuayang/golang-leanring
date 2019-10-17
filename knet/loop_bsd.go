// +build darwin netbsd freebsd openbsd dragonfly

package knet

import (
	"knet/netpoll"
)

func (l *loop) eventHandler(fd int, filter int16) error {
	if fd == 0 {
		return nil
	}
	if c, ok := l.conns[fd]; ok {
		switch {
		case !c.outBuffer.IsEmpty():
			if filter == netpoll.EVFilterWrite {
				return l.loopOut(c)
			}
		case filter == netpoll.EVFilterRead:
			return l.loopIn(c)
		case filter == netpoll.EVFilterSock:
			return l.loopCloseConn(c)
		}
	} else {
		return l.loopAccept(fd)
	}
	return nil
}
