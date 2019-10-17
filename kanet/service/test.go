package service

import (
	"kanet"
	"kanet/actor"
	"kanet/base/timewheel"
	"time"
)

type TestService struct {
	Base
	timewheel *timewheel.TimeWheel
	ch        chan bool
}

type ITestService interface {
	// actor.IActor
	IBase
}

func (test *TestService) SetCh(ch chan bool) {
	test.ch = ch
}

func (test *TestService) Init(chanNum int) {
	test.Actor.Init(chanNum)
	lastTickTime := time.Now().UnixNano() / 1000000
	// count := 1
	var onTick func()
	onTick = func() {
		n := time.Now().UnixNano() / 1000000
		dlt := n - lastTickTime
		println("test on tick:", dlt, "-", test.ActorID())
		lastTickTime = n
		// test.TimeOut(kanet.MilliSecond(1000), onTick)
	}
	// test.RegisterTimer(kanet.MilliSecond(30))
	test.RegisterCall("test_cb", func(n int, s string, m map[int]string, sl []int) {
		// fmt.Printf("[%d] revice from [%d] - [%d]\n", test.ActorID(), test.CallID(), count)
		// fmt.Println("n:", n)
		// fmt.Println("s:", s)
		// fmt.Println("m:", m)
		// fmt.Println("sl:", sl)
		n++
		s = s + " lee"
		m[999] = "adsakldalskdla"
		sl[0] = 23214214
		test.ch <- true
	})

	test.TimeTick(kanet.MilliSecond(30), onTick)
	actor.MGR.AddActor(test)
	test.Start()
}
