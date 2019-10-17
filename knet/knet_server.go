package knet

import (
	"knet/netpoll"
	"sync"
)

type server struct {
	mainLoop     *loop
	subLoopGroup IEventLoopGroup
	wg           sync.WaitGroup
	cond         *sync.Cond
	once         sync.Once
	listeners    map[int]*listener
}

func (svr *server) waitForShutdown() {
	svr.cond.L.Lock()
	svr.cond.Wait()
	svr.cond.L.Unlock()
}

func (svr *server) signalShutdown() {
	svr.once.Do(func() {
		svr.cond.L.Lock()
		svr.cond.Signal()
		svr.cond.L.Unlock()
	})
}

func (svr *server) startLoops() {
	svr.subLoopGroup.iterate(func(i int, l *loop) bool {
		svr.wg.Add(1)
		go func() {
			// l.loopRun()
			svr.activateSubReactor(l)
			svr.wg.Done()
		}()
		return true
	})
}

func (svr *server) closeLoops() {
	svr.subLoopGroup.iterate(func(i int, l *loop) bool {
		l.poller.Close()
		return true
	})
}

func (svr *server) activeLoops(loopNum int) error {
	for i := 0; i < loopNum; i++ {
		if p, err := netpoll.OpenPoller(); err == nil {
			l := &loop{
				idx:    i,
				poller: p,
				conns:  make(map[int]*connection),
				svr:    svr,
				packet: make([]byte, 0xFFFF),
			}
			// l.poller.AddRead()
			svr.subLoopGroup.register(l)
		} else {
			return err
		}
	}
	svr.startLoops()

	if p, err := netpoll.OpenPoller(); err == nil {
		l := &loop{
			idx:    -1,
			poller: p,
			svr:    svr,
			packet: make([]byte, 0xFFFF),
		}
		svr.mainLoop = l
		svr.wg.Add(1)
		go func() {
			svr.activateMainReactor()
			svr.wg.Done()
		}()
	} else {
		return err
	}

	return nil
}

func (svr *server) start(loopNum int) error {
	return svr.activeLoops(loopNum)
}

func (svr *server) stop() {
	// svr.waitForShutdown()
	svr.subLoopGroup.iterate(func(i int, l *loop) bool {
		printIfError(l.poller.Trigger(func() error {
			return ErrClosing
		}))
		return true
	})

	printIfError(svr.mainLoop.poller.Trigger(func() error {
		return ErrClosing
	}))

	svr.wg.Wait()

	svr.subLoopGroup.iterate(func(i int, l *loop) bool {
		for fd, c := range l.conns {
			printIfError(l.loopCloseConn(c))
			delete(l.conns, fd)
		}
		return true
	})

	svr.closeLoops()
	printIfError(svr.mainLoop.poller.Close())
}

func serve(loopNum int) error {
	svr := new(server)

	svr.listeners = make(map[int]*listener)
	svr.subLoopGroup = new(eventLoopGroup)
	svr.cond = sync.NewCond(&sync.Mutex{})

	err := svr.start(loopNum)
	if err != nil {
		svr.closeLoops()
		return err
	}
	// defer svr.stop()

	inst = svr

	return nil
}

func (svr *server) listen(ln *listener, ev eventHandler) error {
	svr.listeners[ln.fd] = ln
	return svr.mainLoop.poller.AddRead(ln.fd)
}
