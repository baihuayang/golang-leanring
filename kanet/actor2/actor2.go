package actor2

import (
	"kanet/base"
	"kanet/base/timewheel"
	"time"
)

const (
	MsgTypeNormal int = iota
	MsgTypeRet
)

const (
	EventDestroy = iota
)

type (
	Actor struct {
		ID          int64
		callFuncMap map[string]interface{}
		ticker      *time.Ticker
		timer       *timewheel.TimeWheel
		queue       *base.MsgQueue
		bStart      bool

		eventChannel chan int
		waitChannel  chan bool
	}

	IActor interface {
		Start()
		Stop()
		Init()
		RegisterTimer()
		RegisterName()
		TimeOut(time.Duration, func())
		TimeTick(time.Duration, func())
	}
)

func (actor *Actor) Start() {
	if !actor.bStart {
		go actor.run()
		actor.bStart = true
	}
}

func (actor *Actor) Stop() {
	actor.eventChannel <- EventDestroy
}

func (actor *Actor) Init() {
	actor.callFuncMap = make(map[string]interface{})
	actor.ticker = time.NewTicker(1<<63 - 1)
	actor.eventChannel = make(chan int)
}

func (actor *Actor) ActorID() int64 {
	return actor.ID
}

func (actor *Actor) RegisterTimer(duration time.Duration) {
	actor.ticker.Stop()
	actor.timer = timewheel.New(duration, 10)
	actor.ticker = time.NewTicker(duration)
}

func (actor *Actor) RegisterName(name string) {
	// MGR.RegisterName(name, actor.actorID)
}

func (actor *Actor) TimeOut(duration time.Duration, cb func()) {
	if actor.timer != nil {
		actor.timer.AddOnceTimer(duration, cb)
	}
}

func (actor *Actor) TimeTick(duration time.Duration, cb func()) {
	if actor.timer != nil {
		actor.timer.AddTimer(duration, cb)
	}
}

//-----------------------------------------

func (actor *Actor) run() {
	for {
		if actor.loop() {
			break
		}
	}
	actor.clear()
}

func (actor *Actor) loop() bool {
	select {
	case event := <-actor.eventChannel:
		if event == EventDestroy {
			return true
		}
	case <-actor.ticker.C:
		if actor.timer != nil {
			actor.timer.TickHandler()
		}
	case <-actor.waitChannel:
		actor.handleMsg()
	}
	return false
}

func (actor *Actor) clear() {
	actor.ID = 0
	actor.bStart = false
	close(actor.eventChannel)
	if actor.ticker != nil {
		actor.ticker.Stop()
	}
	for i := range actor.callFuncMap {
		delete(actor.callFuncMap, i)
	}
}

func (actor *Actor) handleMsg() {

}
