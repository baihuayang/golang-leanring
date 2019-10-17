// Copyright 2019 Andy Pan. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build darwin netbsd freebsd openbsd dragonfly

package knet

import "knet/netpoll"

func (svr *server) activateMainReactor() {
	defer svr.signalShutdown()

	_ = svr.mainLoop.poller.Polling(func(fd int, filter int16) error {
		return svr.acceptNewConnection(fd)
	})
}

func (svr *server) activateSubReactor(lp *loop) {
	defer svr.signalShutdown()

	_ = lp.poller.Polling(func(fd int, filter int16) error {
		c := lp.conns[fd]
		switch {
		case !c.outBuffer.IsEmpty():
			if filter == netpoll.EVFilterWrite {
				return lp.loopOut(c)
			}
		case filter == netpoll.EVFilterRead:
			return lp.loopIn(c)
		case filter == netpoll.EVFilterSock:
			return lp.loopCloseConn(c)
		}
		return nil
	})
}
