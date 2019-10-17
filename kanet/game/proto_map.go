package game

import (
	pb "kanet/proto/pbgo"
	"reflect"
	"strconv"
)

type (
	Proto struct {
		ProtoID  int
		Type     int //0-C2S !0-S2C
		Request  reflect.Type
		ResPonse reflect.Type
		Desc     string
		Module   string
		Key      string
	}
	newproto map[string]*Proto
)

var Protos map[int]*Proto

func addProto(m map[string]*Proto) {
	for key, p := range m {
		if _, ok := Protos[p.ProtoID]; ok {
			panic("same proto id :" + strconv.FormatInt(int64(p.ProtoID), 16))
		}
		p.Key = key
		Protos[p.ProtoID] = p
	}
}

func init() {
	// print("proto map init")
	Protos = make(map[int]*Proto)
	addProto(MTest)
}

func typeof(v interface{}) reflect.Type {
	return reflect.TypeOf(v)
}

var (
	MTest = newproto{
		"test": &Proto{
			ProtoID:  0x0001,
			Request:  typeof(pb.Data{}),
			ResPonse: typeof(pb.Data{}),
		},
	}
)

func InitProtoMap() {
	addProto(MTest)
}
