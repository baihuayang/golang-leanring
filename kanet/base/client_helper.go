package base

import (
	"errors"
	"kanet/game"
	"log"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var ErrNotMatchProto = errors.New("not match protoID")

const (
	PACKETHEADSIZE  int = 2
	PACKETERRNOSIZE     = 2
)

func ClientPack(protoID int, data interface{}, errno int) (buf []byte, err error) {
	pbuf, err := proto.Marshal(data.(proto.Message))
	if err != nil {
		return nil, err
	}
	bufSize := len(pbuf) + PACKETERRNOSIZE + PackHeadSize + 1
	buf = make([]byte, 0, bufSize)
	buf = append(buf, Int16ToBytes(int16(bufSize-PackHeadSize))...)
	buf = append(buf, Int16ToBytes(int16(errno))...)
	buf = append(buf, Int16ToBytes(int16(protoID))...)
	buf = append(buf, 0)
	buf = append(buf, pbuf...)
	// append(buf, )
	return
}

func ClientUnPack(protoID int, buf []byte) (data interface{}, err error) {
	p, ok := game.Protos[protoID]
	if !ok {
		log.Println("no proto id:", protoID)
		err = ErrNotMatchProto
		return
	}

	uncompress := buf[0]
	if uncompress != 0 {

	}

	val := reflect.New(p.Request)
	packet := val.Interface().(proto.Message)
	err = proto.Unmarshal(buf[1:], packet)
	if err != nil {
		log.Println("client unpack err:", err)
		return
	}
	data = val.Elem().Interface()

	return
}
