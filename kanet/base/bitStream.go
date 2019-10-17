package base

//----------------bitsream---------------
//for example
//buf := make([]byte, 256)
//var bitstream base.BitStream
//bitstream.BuildPacketStream(buf, 256)
//bitstream.WriteInt(1000, 16)
// or
//bitstream := NewBitStream(buf)
//----------------------------------------

const (
	Bit8         = 8
	Bit16        = 16
	Bit32        = 32
	Bit64        = 64
	Bit128       = 128
	MaxPacket    = 128 * 1024
	PackHeadSize = 4
)

type (
	BitStream struct {
		dataPtr        []byte
		bitNum         int
		flagNum        int
		tailFlag       bool
		bufSzie        int
		bitsLimit      int
		error          bool
		maxReadBitNum  int
		maxWriteBitNum int
		zipFlag        bool
		tmpBuf         []byte //zip buf
		zipSize        int    //zip length
	}

	IBitStream interface {
		BuildPacketStream([]byte, int) bool
		setBuffer([]byte, int, int)
		GetBuffer() []byte
		GetBytePtr() []byte
		GetReadByteSize() int
		GetCurPos() int
		GetPosition() int
		GetStreamSize() int
		SetPosition(int) bool
		clear()
		resize() bool

		WriteBits(int, []byte)
		ReadBits(int) []byte
		WriteInt(int)
		ReadInt() int
		ReadFlag() bool
		WriteFlag(bool) bool
		ReadInt16() int16
		WriteInt16(int16)
		ReadInt32() int32
		WriteInt32(int32)
		WriteString(string)
		ReadString() string
		WriteInt8(int8)
		ReadInt8() int8
		WriteInt64(int64, int)
		ReadInt64(int) int64
		WriteFloat(float32)
		ReadFloat() float32
		WriteFloat64(float64)
		ReadFloat64() float64
	}
)

func (bitstream *BitStream) setBuffer(bufPtr []byte, size int, maxWriteSize int) {
	bitstream.dataPtr = bufPtr
	bitstream.bitNum = 0
	bitstream.flagNum = 0
	bitstream.tailFlag = false
	bitstream.bufSzie = size
	bitstream.maxReadBitNum = size << 3
	if maxWriteSize < 0 {
		maxWriteSize = size
	}
	bitstream.maxWriteBitNum = maxWriteSize << 3
	bitstream.bitsLimit = size
	bitstream.error = false
}

func (bitstream *BitStream) GetBuffer() []byte {
	return bitstream.dataPtr[0:bitstream.GetPosition()]
}

func (bitstream *BitStream) GetBytePtr() []byte {
	return bitstream.dataPtr[bitstream.GetPosition():]
}

func (bitstream *BitStream) GetPosition() int {
	return (bitstream.bitNum + 7) >> 3
}

func (bitstream *BitStream) SetPosition(pos int) bool {
	Assert(pos == 0 || bitstream.flagNum == 0, "不正确的SetPosition调用")
	if pos != 0 && bitstream.flagNum != 0 {
		return false
	}
	bitstream.bitNum = pos << 3
	bitstream.flagNum = 0
	return true
}

func (bitstream *BitStream) BuildPacketSteam(buffer []byte, writeSize int) bool {
	if writeSize <= 0 {
		return false
	}
	bitstream.setBuffer(buffer, writeSize, -1)
	bitstream.SetPosition(0)
	return true
}

func (bitstream *BitStream) GetReadByteSize() int {
	return (bitstream.maxReadBitNum >> 3) - bitstream.GetPosition()
}

func (bitstream *BitStream) GetCurPos() int {
	return bitstream.bitNum
}

func (bitstream *BitStream) GetStreamSize() int {
	return bitstream.bufSzie
}

func (bitstream *BitStream) clear() {
	// var buff []byte
	buff := make([]byte, bitstream.bufSzie)
	bitstream.dataPtr = buff
}

func (bitstream *BitStream) resize() bool {
	bitstream.dataPtr = append(bitstream.dataPtr, make([]byte, bitstream.bitsLimit)...)
	size := bitstream.bitsLimit * 2

	bitstream.bitsLimit = size
	bitstream.bufSzie = size
	bitstream.maxWriteBitNum = size << 3
	bitstream.maxWriteBitNum = size << 3

	return true
}

func (bitstream *BitStream) WriteBits(bitCount int, bitPtr []byte) {
	if bitCount < 0 {
		return
	}

	if (bitCount & 0x7) != 0 { //向下对8取整
		bitCount = (bitCount & ^0x7) + 8
	}

	if bitstream.tailFlag {
		bitstream.error = true
		Assert(false, "out of range write[WriteBits]")
		return
	}

	if bitstream.bitNum+bitCount > bitstream.maxWriteBitNum {
		if !bitstream.resize() {
			bitstream.error = true
			Assert(false, "out of range write[WriteBits]")
			return
		}
	}

	byteNum := bitstream.bitNum >> 3
	byteCount := (bitCount + 7) >> 3 //向下取整??

	for i, v := range bitPtr[:byteCount] {
		bitstream.dataPtr[byteNum+i] = v
	}

	bitstream.bitNum += bitCount
}

func (bitstream *BitStream) ReadBits(bitCount int) []byte {
	if bitCount < 0 {
		return nil
	}

	if (bitCount & 0x7) != 0 { //向下对8取整
		bitCount = (bitCount & ^0x7) + 8
	}

	if bitstream.tailFlag {
		bitstream.error = true
		Assert(false, "out of range read[ReadBits]")
		return nil
	}

	if bitstream.bitNum+bitCount > bitstream.maxWriteBitNum {
		bitstream.error = true
		Assert(false, "out of range read[ReadBits]")
		return nil
	}

	byteNum := bitstream.bitNum >> 3
	byteCount := (bitCount + 7) >> 3

	bitstream.bitNum += bitCount

	return bitstream.dataPtr[byteNum : byteNum+byteCount]
}

func (bitstream *BitStream) WriteInt8(val int8) {
	bitstream.WriteBits(Bit8, Int8ToBytes(val))
}

func (bitstream *BitStream) WriteInt(val int, bitCount int) {
	bitstream.WriteBits(bitCount, IntToBytes(val))
}

func (bitstream *BitStream) ReadInt(bitCount int) int {
	buff := bitstream.ReadBits(bitCount)
	ret := BytesToInt(buff)
	if bitCount == Bit32 {
		return ret
	}
	return ret & ((1 << uint32(bitCount)) - 1) //???
}

func (bitstream *BitStream) WriteInt32(val int32) {
	bitstream.WriteBits(Bit32, Int32ToBytes(val))
}

func (bitstream *BitStream) ReadInt32() int32 {
	buff := bitstream.ReadBits(Bit32)
	ret := BytesToInt32(buff)
	// if bitCount == Bit32 {
	return ret
	// }
	// return ret & ((1 << uint32(bitCount)) - 1) //???
}

func (bitstream *BitStream) WriteInt16(val int16) {
	bitstream.WriteBits(Bit16, Int16ToBytes(val))
}

func (bitstream *BitStream) ReadInt16() int16 {
	buff := bitstream.ReadBits(Bit16)
	ret := BytesToInt16(buff)
	// if bitCount == Bit32 {
	return ret
	// }
	// return ret & ((1 << uint32(bitCount)) - 1) //???
}

func (bitstream *BitStream) WriteFlag(val bool) bool {
	if (bitstream.flagNum-((bitstream.flagNum>>3)<<3) == 0) && !bitstream.tailFlag {
		bitstream.flagNum = bitstream.bitNum
		if bitstream.bitNum+8 < bitstream.maxWriteBitNum {
			bitstream.bitNum += 8
		} else {
			if !bitstream.resize() {
				bitstream.tailFlag = true
			} else {
				bitstream.bitNum += 8
			}
		}
	}

	if bitstream.flagNum+1 >= bitstream.maxWriteBitNum {
		bitstream.error = true
		Assert(false, "Out of range write")
		return false
	}

	if val {
		// xx & 0x7 求余8
		bitstream.dataPtr[bitstream.flagNum>>3] |= (1 << uint32(bitstream.flagNum&0x7))
	} else {
		bitstream.dataPtr[bitstream.flagNum>>3] &= ^(1 << uint32(bitstream.flagNum&0x7))
	}

	bitstream.flagNum++
	return val
}

func (bitstream *BitStream) ReadFlag() bool {
	if (bitstream.flagNum-((bitstream.flagNum>>3)<<3) == 0) && !bitstream.tailFlag {
		bitstream.flagNum = bitstream.bitNum
		if bitstream.bitNum+8 < bitstream.maxWriteBitNum {
			bitstream.bitNum += 8
		} else {
			if !bitstream.resize() {
				bitstream.tailFlag = true
			} else {
				bitstream.bitNum += 8
			}
		}
	}

	if bitstream.flagNum+1 >= bitstream.maxReadBitNum {
		bitstream.error = true
		Assert(false, "Out of range read")
		return false
	}

	ret := int(bitstream.dataPtr[bitstream.flagNum>>3]&(1<<uint32(bitstream.flagNum&0x7))) != 0

	bitstream.flagNum++

	return ret
}

func (bitstream *BitStream) ReadString() string {
	if bitstream.ReadFlag() {
		nLen := bitstream.ReadInt16()
		buf := bitstream.ReadBits(int(nLen) << 3)
		return string(buf)
	}
	return string("")
}

func (bitstream *BitStream) WriteString(str string) {
	buf := []byte(str)
	nLen := len(buf)
	if bitstream.WriteFlag(nLen > 0) {
		bitstream.WriteInt16(int16(nLen))
		bitstream.WriteBits(nLen<<3, buf)
	}
}

func (bitstream *BitStream) WriteInt64(val int64) {
	bitstream.WriteBits(Bit64, Int64ToBytes(val))
}

func (bitstream *BitStream) ReadInt64() int64 {
	buff := bitstream.ReadBits(Bit64)
	ret := BytesToInt64(buff)
	// if bitCount == Bit64 {
	return ret
	// }
	// return ret & ((1 << uint64(bitCount)) - 1) //???
}

func (bitstream *BitStream) WriteFloat(value float32) {
	bitstream.WriteBits(Bit32, Float32ToBytes(value))
}

func (bitstream *BitStream) ReadFloat() float32 {
	buf := bitstream.ReadBits(Bit32)
	return BytesToFloat32(buf)
}

func (bitstream *BitStream) WriteFloat64(value float64) {
	bitstream.WriteBits(Bit64, Float64ToBytes(value))
}

func (bitstream *BitStream) ReadFloat64() float64 {
	buf := bitstream.ReadBits(Bit64)
	return BytesToFloat64(buf)
}

func NewBitStream(buf []byte, nLen int) *BitStream {
	var bitstream BitStream
	bitstream.BuildPacketSteam(buf, nLen)
	return &bitstream
}
