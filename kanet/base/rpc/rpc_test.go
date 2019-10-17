package rpc

import (
	"kanet/base"
	"reflect"
	"testing"
)

type (
	TopRank struct {
		Value []int `sql:"name:value" json:"value"	json:"value"`
	}
)

var (
	ntimes     = 100000
	nArraySize = 10
	nValue     = 0x7fffffff
)

func TestRpc(t *testing.T) {
	aa := []int32{}
	for i := 0; i < nArraySize; i++ {
		aa = append(aa, int32(nValue))
	}
	for i := 0; i < ntimes; i++ {
		Pack("test", aa)
	}
}

func test(aa []int32) {

}

func TestURpc(t *testing.T) {
	aa := []int32{}
	for i := 0; i < nArraySize; i++ {
		aa = append(aa, int32(nValue))
	}
	buff := Pack("test", aa)
	funcType := reflect.TypeOf(test)
	for i := 0; i < ntimes; i++ {
		bitstream := base.NewBitStream(buff, len(buff))
		UnPack(funcType, bitstream)
	}
}
