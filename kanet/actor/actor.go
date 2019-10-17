package actor

import (
	"fmt"
	"kanet/base"
	"kanet/base/rpc"
	"kanet/base/timewheel"
	"log"
	"reflect"
	"strings"
	"time"
)

const (
	MsgTypeNormal = iota
	MsgTypeRet
	MsgTypeClient
)

type (
	Actor struct {
		msgChannel    chan CallIO    //rpc chan
		localChannel  chan LocalCall //local chan
		statusChannel chan int       //use for status
		// timerChannel  chan bool
		actorID     int64
		callFuncMap map[string]interface{}
		asynCallMap map[int]interface{}
		ticker      *time.Ticker
		timer       *timewheel.TimeWheel

		start    bool
		socketID int
		callID   int64
		onTick   func()

		session int
	}

	IActor interface {
		Init(chanNum int)
		Stop()
		Start()
		FindCall(funcName string) interface{}
		RegisterCall(funcName string, call interface{})
		RegisterName(string)
		SendMsg(interface{}, string, ...interface{})
		SendByName(string, string, ...interface{})
		SendByID(actorID int64, funcName string, args ...interface{})
		CallMsg(interface{}, interface{}, string, ...interface{})
		CallByName(interface{}, string, string, ...interface{})
		CallByID(interface{}, int64, string, ...interface{})
		Send(io CallIO)
		LocalSend(lc LocalCall)
		PacketFunc(id int, buff []byte) bool //回调函数
		RegisterTimer(duration time.Duration)
		ActorID() int64
		CallID() int64
		SocketID() int
		NextSession() int
		SendNoBlock(io CallIO)
		// Tick()
		TimeOut(time.Duration, func())
		TimeTick(time.Duration, func())

		Callback(io CallIO)
		LocalCb(lc LocalCall)
		// SetOnTick(func())
		// SetTimer(time.Duration)
	}

	CallIO struct {
		SocketID int
		Source   int64
		Buff     []byte
		Type     uint8
		Session  int
	}

	LocalCall struct {
		Fname   string
		Args    []interface{}
		Source  int64
		Type    uint8
		Session int
	}
)

const (
	DestroyEvent = iota
)

func (actor *Actor) ActorID() int64 {
	return actor.actorID
}

func (actor *Actor) SocketID() int {
	return actor.socketID
}

func (actor *Actor) NextSession() int {
	actor.session++
	if actor.session >= (1<<31)-1 {
		actor.session = 1
	}
	return actor.session
}

func (actor *Actor) CallID() int64 {
	return actor.callID
}

func (actor *Actor) Init(chanNum int) {
	actor.msgChannel = make(chan CallIO, chanNum)
	actor.localChannel = make(chan LocalCall, chanNum)
	actor.statusChannel = make(chan int, 1)
	actor.actorID = AssignActorID()
	actor.callFuncMap = make(map[string]interface{})
	actor.asynCallMap = make(map[int]interface{})
	actor.ticker = time.NewTicker(1<<63 - 1)
	// actor.timer = time.NewTicker(1<<63 - 1)
	// actor.timerChannel = make(chan bool)
	// actor.onTick = nil
}

func (actor *Actor) RegisterName(name string) {
	MGR.RegisterName(name, actor.actorID)
}

func (actor *Actor) RegisterTimer(duration time.Duration) {
	actor.ticker.Stop()
	actor.timer = timewheel.New(duration, 12)
	actor.ticker = time.NewTicker(duration)
	// actor.timer.Start()
	// actor.timer.Stop()
	// actor.timer = time.NewTicker(duration)
	// actor.onTick = cb
}

// func (actor *Actor) RegisterTimer(duration time.Duration, fun interface{}) {
// 	// actor.timer.Stop()
// 	// actor.timer = time.NewTicker(duration)
// 	actor.onTick = fun.(func())
// }

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

// func (actor *Actor) SetOnTick(onTick func()) {
// 	actor.onTick = onTick
// }

func (actor *Actor) clear() {
	actor.actorID = 0
	actor.callID = 0
	actor.socketID = 0
	actor.start = false
	close(actor.msgChannel)
	close(actor.statusChannel)
	if actor.ticker != nil {
		actor.ticker.Stop()
	}
	for i := range actor.callFuncMap {
		delete(actor.callFuncMap, i)
	}

	for i := range actor.asynCallMap {
		delete(actor.asynCallMap, i)
	}
}

func (actor *Actor) Stop() {
	actor.statusChannel <- DestroyEvent
}

func (actor *Actor) Start() {
	if actor.start == false {
		go actor.run()
		actor.start = true
	}
}

// func (actor *Actor) Tick() {
// 	actor.timerChannel <- true
// }

func (actor *Actor) loop() bool {
	defer func() {
		if err := recover(); err != nil {
			// base.TraceCode(err)
			log.Println(err)
		}
	}()
	select {
	case io := <-actor.msgChannel:
		actor.Callback(io)
	case lc := <-actor.localChannel:
		actor.LocalCb(lc)
	case msg := <-actor.statusChannel:
		if msg == DestroyEvent {
			return true
		}
	// case <-actor.timerChannel:
	// 	actor.onTick()
	case <-actor.ticker.C:
		if actor.timer != nil {
			actor.timer.TickHandler()
		}
	}
	return false
}

func (actor *Actor) run() {
	for {
		if actor.loop() {
			break
		}
	}
	actor.clear()
}

func (actor *Actor) FindCall(funcName string) interface{} {
	funcName = strings.ToLower(funcName)
	fun, exist := actor.callFuncMap[funcName]
	if exist == true {
		return fun
	}
	return nil
}

func (actor *Actor) RegisterCall(funcName string, call interface{}) {
	switch call.(type) {
	case func(*IActor, []byte):
		log.Fatalln("actor error 消息定义函数不符合:", funcName)
		return
	}
	funcName = strings.ToLower(funcName)
	if actor.FindCall(funcName) != nil {
		log.Fatalln("actor error 消息重复定义:", funcName)
		return
	}
	actor.callFuncMap[funcName] = call
}

func (actor *Actor) SendMsg(key interface{}, funcName string, params ...interface{}) {
	MGR.Send(actor.actorID, 0, key, funcName, params...)
}

func (actor *Actor) SendByName(name string, funcName string, args ...interface{}) {
	MGR.SendByName(actor.actorID, 0, name, funcName, args...)
}

func (actor *Actor) SendByID(actorID int64, funcName string, args ...interface{}) {
	MGR.SendByID(actor.actorID, 0, actorID, funcName, args...)

}

func (actor *Actor) CallMsg(cb interface{}, key interface{}, funcName string, params ...interface{}) {
	session := actor.NextSession()
	actor.asynCallMap[session] = cb
	MGR.Send(actor.actorID, session, key, funcName, params...)
}

func (actor *Actor) CallByName(cb interface{}, name string, funcName string, args ...interface{}) {
	session := actor.NextSession()
	actor.asynCallMap[session] = cb
	MGR.SendByName(actor.actorID, session, name, funcName, args...)
}

func (actor *Actor) CallByID(cb interface{}, actorID int64, funcName string, args ...interface{}) {
	session := actor.NextSession()
	actor.asynCallMap[session] = cb
	MGR.SendByID(actor.actorID, session, actorID, funcName, args...)

}

func (actor *Actor) LocalSend(lc LocalCall) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("local Send", err)
		}
	}()
	actor.localChannel <- lc
}

func (actor *Actor) Send(io CallIO) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Send", err)
		}
	}()
	actor.msgChannel <- io
}

func (actor *Actor) SendNoBlock(io CallIO) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("SendNoBlock", err)
		}
	}()
	select {
	case actor.msgChannel <- io:
	default:
		break
	}
}

func (actor *Actor) PacketFunc(id int, buff []byte) bool {
	var io CallIO
	io.Buff = buff
	io.SocketID = id
	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	// funcName = strings.ToLower(funcName)
	pFunc := actor.FindCall(funcName)
	if pFunc != nil {
		actor.Send(io)
		return true
	}
	return false
}

func (actor *Actor) Callback(io CallIO) {
	if f := actor.callFuncMap["__filter"]; f != nil {
		if f.(func(io CallIO) bool)(io) {
			return
		}
	}
	// if io.Type == MsgTypeClient {
	// 	actor.redirect(io.Buff)
	// 	return
	// }
	bitstream := base.NewBitStream(io.Buff, len(io.Buff))
	funcName := bitstream.ReadString()
	// funcName = strings.ToLower(funcName)
	// println("funcName", funcName)
	pFunc := actor.FindCall(funcName)
	if pFunc == nil {
		return
	}

	f := reflect.ValueOf(pFunc)
	k := reflect.TypeOf(pFunc)

	args := rpc.UnPack(k, bitstream)

	if k.NumIn() != len(args) {
		log.Printf("func [%s] can not call , args [%v]", funcName, args)
		return
	}
	actor.callID = io.Source
	if len(args) >= 1 {
		in := make([]reflect.Value, len(args))
		for i, param := range args {
			in[i] = reflect.ValueOf(param)
			if k.In(i).Kind() != in[i].Kind() {
				log.Printf("func [%s] args no fit, args [func(%v)]", funcName, in)
				return
			}
		}
		f.Call(in)
	} else {
		f.Call(nil)
	}
}

func (actor *Actor) LocalCb(lc LocalCall) {
	// if lc.Type == MsgTypeClient {
	// 	actor.redirect(lc.Args[0].([]byte))
	// 	return
	// }
	if f := actor.callFuncMap["__filter"]; f != nil {
		if f.(func(lc LocalCall) bool)(lc) {
			return
		}
	}
	var pFunc interface{}
	if lc.Type == MsgTypeRet && lc.Session > 0 {
		c, ok := actor.asynCallMap[lc.Session]
		if !ok {
			return
		}
		pFunc = c
		delete(actor.asynCallMap, lc.Session)
	} else {
		pFunc = actor.FindCall(lc.Fname)
	}

	if pFunc == nil {
		return
	}

	fv := reflect.ValueOf(pFunc)
	ft := reflect.TypeOf(pFunc)

	if ft.NumIn() != len(lc.Args) {
		log.Printf("func [%s] can not call , arg [%v]", lc.Fname, lc.Args)
		return
	}
	if len(lc.Args) > 1 {
		in := make([]reflect.Value, len(lc.Args))
		for i, arg := range lc.Args {
			in[i] = reflect.ValueOf(arg)
			if ft.In(i).Kind() != in[i].Kind() {
				log.Printf("func [%s] arg no fit, arg [func(%v)]", lc.Fname, in)
				return
			}
		}
		if lc.Type != MsgTypeRet && lc.Session > 0 {
			actor.ret(lc.Source, lc.Session, fv.Call(in))
		} else {
			fv.Call(in)
		}
	} else {
		if lc.Type != MsgTypeRet && lc.Session > 0 {
			actor.ret(lc.Source, lc.Session, fv.Call(nil))
		} else {
			fv.Call(nil)
		}
	}
}

func (actor *Actor) ret(source int64, session int, values []reflect.Value) {
}
