package service

import (
	"kanet/actor"
	"kanet/base"
	"kanet/game"
	"log"
)

type (
	Base struct {
		actor.Actor
		modules map[string]interface{}
	}
	IBase interface {
		actor.IActor
	}
)

func (b *Base) redirect(buf []byte) {
	protoID := base.BytesToInt(buf[:2])
	itf, err := base.ClientUnPack(protoID, buf[2:])
	if err != nil {
		log.Println(err)
		return
	}
	p := game.Protos[protoID]
	if m, ok := b.modules[p.Module]; ok {
		println(itf, m)
	}
}

func (b *Base) filter(io actor.CallIO) bool {
	if io.Type == actor.MsgTypeClient {
		b.redirect(io.Buff)
		return true
	}
	return false
}

func (b *Base) localFilter(lc actor.LocalCall) bool {
	if lc.Type == actor.MsgTypeClient {
		b.redirect(lc.Args[0].([]byte))
		return true
	}
	return false
}

func (b *Base) Init(chanNum int) {
	b.RegisterCall("__filter", b.filter)
	b.RegisterCall("__localfilter", b.localFilter)
}

// func (b *Base) Callback(io actor.CallIO) {
// 	if io.Type == actor.MsgTypeClient {
// 		b.redirect(io.Buff)
// 		return
// 	}
// 	b.Actor.Callback(io)
// }

// func (b *Base) LocalCb(lc actor.LocalCall) {
// 	println("service base local callback")
// 	if lc.Type == actor.MsgTypeClient {
// 		b.redirect(lc.Args[0].([]byte))
// 		return
// 	}
// 	b.Actor.LocalCb(lc)
// }
