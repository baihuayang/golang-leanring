package actor

import (
	"kanet/lib/snowflake"
	"log"
	"reflect"
	"runtime"
)

var (
	gIDSeed int64
	MGR     *actorMgr
)

type actorMgr struct {
	actorIDMap     map[int64]IActor
	actorNameIDMap map[string]int64
}

func AssignActorID() int64 {
	// atomic.AddInt64(&gIDSeed, 1)
	// return int64(gIDSeed)
	return snowflake.UUID()
}

func (mgr *actorMgr) NewActor(chanNum int) int64 {
	pActor := new(Actor)
	pActor.Init(chanNum)
	mgr.AddActor(pActor)
	return pActor.ActorID()
}

func (mgr *actorMgr) Init() {
	mgr.actorIDMap = make(map[int64]IActor)
	mgr.actorNameIDMap = make(map[string]int64)
}

func (mgr *actorMgr) AddActor(actor IActor) {
	mgr.actorIDMap[actor.ActorID()] = actor
}

func (mgr *actorMgr) RegisterName(name string, actorID int64) {
	_, exist := mgr.actorIDMap[actorID]
	if exist {
		mgr.actorNameIDMap[name] = actorID
	}
}

func (mgr *actorMgr) GetActor(actorID int64) IActor {
	pActor, exist := mgr.actorIDMap[actorID]
	if exist {
		return pActor
	}
	return nil
}

func (mgr *actorMgr) send(id int64, funcName string, args ...interface{}) {

}

func (mgr *actorMgr) SendByName(source int64, session int, name string, funcName string, args ...interface{}) {
	ID, exist := mgr.actorNameIDMap[name]
	if exist {
		mgr.SendByID(source, session, ID, funcName, args...)
	}

}

func (mgr *actorMgr) SendByID(source int64, session int, actorID int64, funcName string, args ...interface{}) {
	iActor, exist := mgr.actorIDMap[actorID]
	if exist {
		lc := LocalCall{
			Fname:  funcName,
			Args:   args,
			Source: source,
		}
		iActor.LocalSend(lc)
	}
}

func (mgr *actorMgr) Send(source int64, session int, key interface{}, funcName string, args ...interface{}) {
	switch key.(type) {
	case int:
		mgr.SendByID(source, session, int64(key.(int)), funcName, args...)
	case int64:
		mgr.SendByID(source, session, key.(int64), funcName, args...)
	case string:
		ID, exist := mgr.actorNameIDMap[funcName]
		if exist {
			mgr.SendByID(source, session, ID, funcName, args...)
		}
	default:
		_, file, line, _ := runtime.Caller(0)
		log.Printf("[%s:%d] wrong key : [%s], kind:[%s]", file, line, key, reflect.TypeOf(key).String())
	}
}

func Init() {
	MGR = new(actorMgr)
	MGR.Init()
}
