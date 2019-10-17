package base

import "testing"

func TestBitStream(t *testing.T) {
	b := Int8ToBytes(42)
	println(BytesToInt8(b))
}
