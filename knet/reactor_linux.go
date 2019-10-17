// Copyright 2019 Andy Pan. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build linux

package gnet

import "kanet/network/knet/netpoll"

func (svr *server) activateMainReactor() {
	defer svr.signalShutdown()

	_ = svr.mainLoop.poller.Polling(func(fd int, ev uint32) error {
		return svr.acceptNewConnection(fd)
	})
}

func (svr *server) activateSubReactor(lp *loop) {
	defer svr.signalShutdown()

	_ = lp.poller.Polling(func(fd int, ev uint32) error {
		c := lp.conns[fd]
		switch {
		case !c.outBuffer.IsEmpty():
			if ev&netpoll.OutEvents != 0 {
				return lp.loopOut(c)
			}
		case ev&netpoll.InEvents != 0:
			return lp.loopIn(c)
		}
		return nil
	})
}
