package rpc

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"kanet/base"
	"log"
	"reflect"
	"unsafe"
)

const (
	RPCInt64   = 10
	RPCUInt64  = 11
	RPCPInt64  = 70
	RPCPUInt64 = 71
	RPCMESSAGE = 120
	RPCGOB     = 121
)

const (
	SIZEBOOL    = unsafe.Sizeof(bool(false))
	SIZEINT     = unsafe.Sizeof(int(0))
	SIZEINT8    = unsafe.Sizeof(int8(0))
	SIZEINT16   = unsafe.Sizeof(int16(0))
	SIZEINT32   = unsafe.Sizeof(int32(0))
	SIZEINT64   = unsafe.Sizeof(int64(0))
	SIZEUINT    = unsafe.Sizeof(uint(0))
	SIZEUINT8   = unsafe.Sizeof(uint8(0))
	SIZEUINT16  = unsafe.Sizeof(uint16(0))
	SIZEUINT32  = unsafe.Sizeof(uint32(0))
	SIZEUINT64  = unsafe.Sizeof(uint64(0))
	SIZEFLOAT32 = unsafe.Sizeof(float32(0))
	SIZEFLOAT64 = unsafe.Sizeof(float64(0))
	SIZESTRING  = unsafe.Sizeof(string(0))
	SIZEPTR     = unsafe.Sizeof(uintptr(0))
) //packet size

func Pack(funcName string, params ...interface{}) []byte {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("pack:", err)
		}
	}()

	msg := make([]byte, 1024)
	bitstream := base.NewBitStream(msg, len(msg))
	bitstream.WriteString(funcName)
	bitstream.WriteInt8(int8(len(params)))

	for _, param := range params {
		sType := base.TypeString(param)
		switch sType {
		case "bool":
			bitstream.WriteInt8(1)
			bitstream.WriteFlag(param.(bool))
		case "float64":
			bitstream.WriteInt8(2)
			bitstream.WriteFloat64(param.(float64))
		case "float32":
			bitstream.WriteInt8(3)
			bitstream.WriteFloat(param.(float32))
		case "int8":
			bitstream.WriteInt8(4)
			bitstream.WriteInt(int(param.(int8)), 8)
		case "uint8":
			bitstream.WriteInt8(5)
			bitstream.WriteInt(int(param.(uint8)), 8)
		case "int16":
			bitstream.WriteInt8(6)
			bitstream.WriteInt(int(param.(int16)), 16)
		case "uint16":
			bitstream.WriteInt8(7)
			bitstream.WriteInt(int(param.(uint16)), 16)
		case "int32":
			bitstream.WriteInt8(8)
			bitstream.WriteInt(int(param.(int32)), 32)
		case "uint32":
			bitstream.WriteInt8(9)
			bitstream.WriteInt(int(param.(uint32)), 32)
		case "int64":
			bitstream.WriteInt8(10)
			bitstream.WriteInt64(param.(int64))
		case "uint64":
			bitstream.WriteInt8(11)
			bitstream.WriteInt64(int64(param.(uint64)))
		case "string":
			bitstream.WriteInt8(12)
			bitstream.WriteString(param.(string))
		case "int":
			bitstream.WriteInt8(13)
			bitstream.WriteInt(param.(int), 32)
		case "uint":
			bitstream.WriteInt8(14)
			bitstream.WriteInt(int(param.(uint)), 32)

		case "[]bool":
			bitstream.WriteInt8(21)
			nLen := len(param.([]bool))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(param.([]bool)[i])
			}
		case "[]float64":
			bitstream.WriteInt8(22)
			nLen := len(param.([]float64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(param.([]float64)[i])
			}
		case "[]float32":
			bitstream.WriteInt8(23)
			nLen := len(param.([]float32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(param.([]float32)[i])
			}
		case "[]int8":
			bitstream.WriteInt8(24)
			nLen := len(param.([]int8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int8)[i]), 8)
			}
		case "[]uint8":
			bitstream.WriteInt8(25)
			nLen := len(param.([]uint8))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint8)[i]), 8)
			}
		case "[]int16":
			bitstream.WriteInt8(26)
			nLen := len(param.([]int16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int16)[i]), 16)
			}
		case "[]uint16":
			bitstream.WriteInt8(27)
			nLen := len(param.([]uint16))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint16)[i]), 16)
			}
		case "[]int32":
			bitstream.WriteInt8(28)
			nLen := len(param.([]int32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]int32)[i]), 32)
			}
		case "[]uint32":
			bitstream.WriteInt8(29)
			nLen := len(param.([]uint32))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint32)[i]), 32)
			}
		case "[]int64":
			bitstream.WriteInt8(30)
			nLen := len(param.([]int64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(param.([]int64)[i])
			}
		case "[]uint64":
			bitstream.WriteInt8(31)
			nLen := len(param.([]uint64))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(param.([]uint64)[i]))
			}
		case "[]string":
			bitstream.WriteInt8(32)
			nLen := len(param.([]string))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(param.([]string)[i])
			}
		case "[]int":
			bitstream.WriteInt8(33)
			nLen := len(param.([]int))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(param.([]int)[i], 32)
			}
		case "[]uint":
			bitstream.WriteInt8(34)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(param.([]uint)[i]), 32)
			}

		case "[*]bool":
			bitstream.WriteInt8(41)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFlag(val.Index(i).Bool())
			}
		case "[*]float64":
			bitstream.WriteInt8(42)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat64(val.Index(i).Float())
			}
		case "[*]float32":
			bitstream.WriteInt8(43)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteFloat(float32(val.Index(i).Float()))
			}
		case "[*]int8":
			bitstream.WriteInt8(44)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 8)
			}
		case "[*]uint8":
			bitstream.WriteInt8(45)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 8)
			}
		case "[*]int16":
			bitstream.WriteInt8(46)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 16)
			}
		case "[*]uint16":
			bitstream.WriteInt8(47)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 16)
			}
		case "[*]int32":
			bitstream.WriteInt8(48)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]uint32":
			bitstream.WriteInt8(49)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}
		case "[*]int64":
			bitstream.WriteInt8(50)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(val.Index(i).Int())
			}
		case "[*]uint64":
			bitstream.WriteInt8(51)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt64(int64(val.Index(i).Uint()))
			}
		case "[*]string":
			bitstream.WriteInt8(52)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteString(val.Index(i).String())
			}
		case "[*]int":
			bitstream.WriteInt8(53)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Int()), 32)
			}
		case "[*]uint":
			bitstream.WriteInt8(54)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				bitstream.WriteInt(int(val.Index(i).Uint()), 32)
			}

		case "*bool":
			if param.(*bool) != nil {
				bitstream.WriteInt8(61)
				bitstream.WriteFlag(*param.(*bool))
			} else {
				bitstream.WriteInt8(61)
				bitstream.WriteFlag(false)
			}
		case "*float64":
			if param.(*float64) != nil {
				bitstream.WriteInt8(62)
				bitstream.WriteFloat64(*param.(*float64))
			} else {
				bitstream.WriteInt8(62)
				bitstream.WriteFloat64(0)
			}
		case "*float32":
			if param.(*float32) != nil {
				bitstream.WriteInt8(63)
				bitstream.WriteFloat(*param.(*float32))
			} else {
				bitstream.WriteInt8(63)
				bitstream.WriteFloat(0)
			}
		case "*int8":
			if param.(*int8) != nil {
				bitstream.WriteInt8(64)
				bitstream.WriteInt(int(*param.(*int8)), 8)
			} else {
				bitstream.WriteInt8(64)
				bitstream.WriteInt8(0)
			}
		case "*uint8":
			if param.(*uint8) != nil {
				bitstream.WriteInt8(65)
				bitstream.WriteInt(int(*param.(*uint8)), 8)
			} else {
				bitstream.WriteInt8(65)
				bitstream.WriteInt8(0)
			}
		case "*int16":
			if param.(*int16) != nil {
				bitstream.WriteInt8(66)
				bitstream.WriteInt(int(*param.(*int16)), 16)
			} else {
				bitstream.WriteInt8(66)
				bitstream.WriteInt(0, 16)
			}
		case "*uint16":
			if param.(*uint16) != nil {
				bitstream.WriteInt8(67)
				bitstream.WriteInt(int(*param.(*uint16)), 16)
			} else {
				bitstream.WriteInt8(67)
				bitstream.WriteInt(0, 16)
			}
		case "*int32":
			if param.(*int32) != nil {
				bitstream.WriteInt8(68)
				bitstream.WriteInt(int(*param.(*int32)), 32)
			} else {
				bitstream.WriteInt8(68)
				bitstream.WriteInt(0, 32)
			}
		case "*uint32":
			if param.(*uint32) != nil {
				bitstream.WriteInt8(69)
				bitstream.WriteInt(int(*param.(*uint32)), 32)
			} else {
				bitstream.WriteInt8(69)
				bitstream.WriteInt(0, 32)
			}
		case "*int64":
			if param.(*int64) != nil {
				bitstream.WriteInt8(70)
				bitstream.WriteInt64(*param.(*int64))
			} else {
				bitstream.WriteInt8(70)
				bitstream.WriteInt64(0)
			}
		case "*uint64":
			if param.(*uint64) != nil {
				bitstream.WriteInt8(71)
				bitstream.WriteInt64(int64(*param.(*uint64)))
			} else {
				bitstream.WriteInt8(71)
				bitstream.WriteInt64(0)
			}
		case "*string":
			if param.(*string) != nil {
				bitstream.WriteInt8(72)
				bitstream.WriteString(*param.(*string))
			} else {
				bitstream.WriteInt8(72)
				bitstream.WriteString("")
			}
		case "*int":
			if param.(*int) != nil {
				bitstream.WriteInt8(73)
				bitstream.WriteInt(*param.(*int), 32)
			} else {
				bitstream.WriteInt8(73)
				bitstream.WriteInt(0, 32)
			}
		case "*uint":
			if param.(*uint) != nil {
				bitstream.WriteInt8(74)
				bitstream.WriteInt(int(*param.(*uint)), 32)
			} else {
				bitstream.WriteInt8(74)
				bitstream.WriteInt(0, 32)
			}

		case "[]*bool":
			bitstream.WriteInt8(81)
			nLen := len(param.([]*bool))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*bool) {
				if v != nil {
					bitstream.WriteFlag(*v)
				} else {
					bitstream.WriteFlag(false)
				}
			}
		case "[]*float64":
			bitstream.WriteInt8(82)
			nLen := len(param.([]float64))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*float64) {
				if v != nil {
					bitstream.WriteFloat64(*v)
				} else {
					bitstream.WriteFloat64(0)
				}
			}
		case "[]*float32":
			bitstream.WriteInt8(83)
			nLen := len(param.([]float32))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*float32) {
				if v != nil {
					bitstream.WriteFloat(*v)
				} else {
					bitstream.WriteFloat(0)
				}
			}
		case "[]*int8":
			bitstream.WriteInt8(84)
			nLen := len(param.([]int8))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int8) {
				if v != nil {
					bitstream.WriteInt(int(*v), 8)
				} else {
					bitstream.WriteInt8(0)
				}
			}
		case "[]*uint8":
			bitstream.WriteInt8(85)
			nLen := len(param.([]uint8))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint8) {
				if v != nil {
					bitstream.WriteInt(int(*v), 8)
				} else {
					bitstream.WriteInt8(0)
				}
			}
		case "[]*int16":
			bitstream.WriteInt8(86)
			nLen := len(param.([]int16))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int16) {
				if v != nil {
					bitstream.WriteInt(int(*v), 16)
				} else {
					bitstream.WriteInt(0, 16)
				}
			}
		case "[]*uint16":
			bitstream.WriteInt8(87)
			nLen := len(param.([]uint16))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint16) {
				if v != nil {
					bitstream.WriteInt(int(*v), 16)
				} else {
					bitstream.WriteInt(0, 16)
				}
			}
		case "[]*int32":
			bitstream.WriteInt8(88)
			nLen := len(param.([]int32))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int32) {
				if v != nil {
					bitstream.WriteInt(int(*v), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*uint32":
			bitstream.WriteInt8(89)
			nLen := len(param.([]uint32))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint32) {
				if v != nil {
					bitstream.WriteInt(int(*v), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}

			}
		case "[]*int64":
			bitstream.WriteInt8(90)
			nLen := len(param.([]int64))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int64) {
				if v != nil {
					bitstream.WriteInt64(*v)
				} else {
					bitstream.WriteInt64(0)
				}
			}
		case "[]*uint64":
			bitstream.WriteInt8(91)
			nLen := len(param.([]uint64))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*uint64) {
				if v != nil {
					bitstream.WriteInt64(int64(*v))
				} else {
					bitstream.WriteInt64(0)
				}
			}
		case "[]*string":
			bitstream.WriteInt8(92)
			nLen := len(param.([]string))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*string) {
				if v != nil {
					bitstream.WriteString(*v)
				} else {
					bitstream.WriteString("")
				}
			}
		case "[]*int":
			bitstream.WriteInt8(93)
			nLen := len(param.([]*int))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int) {
				if v != nil {
					bitstream.WriteInt(*v, 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}
		case "[]*uint":
			bitstream.WriteInt8(94)
			nLen := len(param.([]uint))
			bitstream.WriteInt(nLen, 16)
			for _, v := range param.([]*int) {
				if v != nil {
					bitstream.WriteInt(int(*v), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}

		case "[*]*bool":
			bitstream.WriteInt8(101)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFlag(val.Index(i).Elem().Bool())
				} else {
					bitstream.WriteFlag(false)
				}
			}
		case "[*]*float64":
			bitstream.WriteInt8(102)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFloat64(val.Index(i).Elem().Float())
				} else {
					bitstream.WriteFloat64(0)
				}
			}
		case "[*]*float32":
			bitstream.WriteInt8(103)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteFloat(float32(val.Index(i).Elem().Float()))
				} else {
					bitstream.WriteFloat(0)
				}
			}
		case "[*]*int8":
			bitstream.WriteInt8(104)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 8)
				} else {
					bitstream.WriteInt8(0)
				}
			}
		case "[*]*uint8":
			bitstream.WriteInt8(105)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 8)
				} else {
					bitstream.WriteInt8(0)
				}
			}
		case "[*]*int16":
			bitstream.WriteInt8(106)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 16)
				} else {
					bitstream.WriteInt(0, 16)
				}
			}
		case "[*]*uint16":
			bitstream.WriteInt8(107)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 16)
				} else {
					bitstream.WriteInt(0, 16)
				}
			}
		case "[*]*int32":
			bitstream.WriteInt8(108)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*uint32":
			bitstream.WriteInt8(109)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*int64":
			bitstream.WriteInt8(110)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt64(val.Index(i).Elem().Int())
				} else {
					bitstream.WriteInt64(0)
				}
			}
		case "[*]*uint64":
			bitstream.WriteInt8(111)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt64(int64(val.Index(i).Elem().Uint()))
				} else {
					bitstream.WriteInt64(0)
				}
			}
		case "[*]*string":
			bitstream.WriteInt8(112)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteString(val.Index(i).Elem().String())
				} else {
					bitstream.WriteString("")
				}
			}
		case "[*]*int":
			bitstream.WriteInt8(113)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Int()), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}
		case "[*]*uint":
			bitstream.WriteInt8(114)
			val := reflect.ValueOf(param)
			nLen := val.Len()
			bitstream.WriteInt(nLen, 16)
			for i := 0; i < nLen; i++ {
				if !val.Index(i).IsNil() {
					bitstream.WriteInt(int(val.Index(i).Elem().Uint()), 32)
				} else {
					bitstream.WriteInt(0, 32)
				}
			}

		case "*gob":
			bitstream.WriteInt8(RPCGOB)
			buf := &bytes.Buffer{}
			enc := gob.NewEncoder(buf)
			enc.Encode(param)
			nLen := buf.Len()
			bitstream.WriteInt(nLen, base.Bit32)
			bitstream.WriteBits(nLen<<3, buf.Bytes())
		default:
			fmt.Println("params type not supported", sType, reflect.TypeOf(param))
			panic("params type not supported")
		}
	}

	return bitstream.GetBuffer()
}

func UnPack(funcType reflect.Type, bitstream *base.BitStream) []interface{} {

	// funcName := ""
	// funcName = strings.ToLower(funcName)
	// pFunc := this.FindCall(funcName)
	// if pFunc != nil {
	// f := reflect.ValueOf(pFunc)
	// k := reflect.TypeOf(pFunc)
	// strParams := reflect.TypeOf(pFunc).String()

	nCurLen := bitstream.ReadInt(8)
	params := make([]interface{}, nCurLen)
	for i := 0; i < nCurLen; i++ {
		switch bitstream.ReadInt(8) {
		case 1:
			params[i] = bitstream.ReadFlag()
		case 2:
			params[i] = bitstream.ReadFloat64()
		case 3:
			params[i] = bitstream.ReadFloat()
		case 4:
			params[i] = int8(bitstream.ReadInt(8))
		case 5:
			params[i] = uint8(bitstream.ReadInt(8))
		case 6:
			params[i] = int16(bitstream.ReadInt(16))
		case 7:
			params[i] = uint16(bitstream.ReadInt(16))
		case 8:
			params[i] = int32(bitstream.ReadInt(32))
		case 9:
			params[i] = uint32(bitstream.ReadInt(32))
		case 10:
			params[i] = int64(bitstream.ReadInt64())
		case 11:
			params[i] = uint64(bitstream.ReadInt64())
		case 12:
			params[i] = bitstream.ReadString()
		case 13:
			params[i] = bitstream.ReadInt(32)
		case 14:
			params[i] = uint(bitstream.ReadInt(32))

		case 21:
			nLen := bitstream.ReadInt(16)
			val := make([]bool, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFlag()
			}
			params[i] = val
		case 22:
			nLen := bitstream.ReadInt(16)
			val := make([]float64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFloat64()
			}
			params[i] = val
		case 23:
			nLen := bitstream.ReadInt(16)
			val := make([]float32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadFloat()
			}
			params[i] = val
		case 24:
			nLen := bitstream.ReadInt(16)
			val := make([]int8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int8(bitstream.ReadInt(8))
			}
			params[i] = val
		case 25:
			nLen := bitstream.ReadInt(16)
			val := make([]uint8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint8(bitstream.ReadInt(8))
			}
			params[i] = val
		case 26:
			nLen := bitstream.ReadInt(16)
			val := make([]int16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int16(bitstream.ReadInt(16))
			}
			params[i] = val
		case 27:
			nLen := bitstream.ReadInt(16)
			val := make([]uint16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint16(bitstream.ReadInt(16))
			}
			params[i] = val
		case 28:
			nLen := bitstream.ReadInt(16)
			val := make([]int32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int32(bitstream.ReadInt(32))
			}
			params[i] = val
		case 29:
			nLen := bitstream.ReadInt(16)
			val := make([]uint32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint32(bitstream.ReadInt(32))
			}
			params[i] = val
		case 30:
			nLen := bitstream.ReadInt(16)
			val := make([]int64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = int64(bitstream.ReadInt64())
			}
			params[i] = val
		case 31:
			nLen := bitstream.ReadInt(16)
			val := make([]uint64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint64(bitstream.ReadInt64())
			}
			params[i] = val
		case 32:
			nLen := bitstream.ReadInt(16)
			val := make([]string, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadString()
			}
			params[i] = val
		case 33:
			nLen := bitstream.ReadInt(16)
			val := make([]int, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = bitstream.ReadInt(32)
			}
			params[i] = val
		case 34:
			nLen := bitstream.ReadInt(16)
			val := make([]uint, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = uint(bitstream.ReadInt(32))
			}
			params[i] = val

		case 41:
			nLen := bitstream.ReadInt(16)
			aa := bool(false)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetBool(bitstream.ReadFlag())
			}
			params[i] = val.Interface()
		case 42:
			nLen := bitstream.ReadInt(16)
			aa := float64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetFloat(bitstream.ReadFloat64())
			}
			params[i] = val.Interface()
		case 43:
			nLen := bitstream.ReadInt(16)
			aa := float32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetFloat(float64(bitstream.ReadFloat()))
			}
			params[i] = val.Interface()
		case 44:
			nLen := bitstream.ReadInt(16)
			aa := int8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
			}
			params[i] = val.Interface()
		case 45:
			nLen := bitstream.ReadInt(16)
			aa := uint8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(8)))
			}
			params[i] = val.Interface()
		case 46:
			nLen := bitstream.ReadInt(16)
			aa := int16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
			}
			params[i] = val.Interface()
		case 47:
			nLen := bitstream.ReadInt(16)
			aa := uint16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(16)))
			}
			params[i] = val.Interface()
		case 48:
			nLen := bitstream.ReadInt(16)
			aa := int32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case 49:
			nLen := bitstream.ReadInt(16)
			aa := uint32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case 50:
			nLen := bitstream.ReadInt(16)
			aa := int64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt64()))
			}
			params[i] = val.Interface()
		case 51:
			nLen := bitstream.ReadInt(16)
			aa := uint64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt64()))
			}
			params[i] = val.Interface()
		case 52:
			nLen := bitstream.ReadInt(16)
			aa := string("")
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetString(bitstream.ReadString())
			}
			params[i] = val.Interface()
		case 53:
			nLen := bitstream.ReadInt(16)
			aa := int(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()
		case 54:
			nLen := bitstream.ReadInt(16)
			aa := uint(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(aa))
			val := reflect.New(tVal).Elem()
			for i := 0; i < nLen; i++ {
				val.Index(i).SetInt(int64(bitstream.ReadInt(32)))
			}
			params[i] = val.Interface()

		case 61:
			val := new(bool)
			*val = bitstream.ReadFlag()
			params[i] = val
		case 62:
			val := new(float64)
			*val = bitstream.ReadFloat64()
			params[i] = val
		case 63:
			val := new(float32)
			*val = bitstream.ReadFloat()
			params[i] = val
		case 64:
			val := new(int8)
			*val = int8(bitstream.ReadInt(8))
			params[i] = val
		case 65:
			val := new(uint8)
			*val = uint8(bitstream.ReadInt(8))
			params[i] = val
		case 66:
			val := new(int16)
			*val = int16(bitstream.ReadInt(16))
			params[i] = val
		case 67:
			val := new(uint16)
			*val = uint16(bitstream.ReadInt(16))
			params[i] = val
		case 68:
			val := new(int32)
			*val = int32(bitstream.ReadInt(32))
			params[i] = val
		case 69:
			val := new(uint32)
			*val = uint32(bitstream.ReadInt(32))
			params[i] = val
		case 70:
			val := new(int64)
			*val = int64(bitstream.ReadInt64())
			params[i] = val
		case 71:
			val := new(uint64)
			*val = uint64(bitstream.ReadInt64())
			params[i] = val
		case 72:
			val := new(string)
			*val = bitstream.ReadString()
			params[i] = val
		case 73:
			val := new(int)
			*val = bitstream.ReadInt(32)
			params[i] = val
		case 74:
			val := new(uint)
			*val = uint(bitstream.ReadInt(32))
			params[i] = val

		case 81:
			nLen := bitstream.ReadInt(16)
			val := make([]*bool, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(bool)
				*val[i] = bitstream.ReadFlag()
			}
			params[i] = val
		case 82:
			nLen := bitstream.ReadInt(16)
			val := make([]*float64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(float64)
				*val[i] = bitstream.ReadFloat64()
			}
			params[i] = val
		case 83:
			nLen := bitstream.ReadInt(16)
			val := make([]*float32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(float32)
				*val[i] = bitstream.ReadFloat()
			}
			params[i] = val
		case 84:
			nLen := bitstream.ReadInt(16)
			val := make([]*int8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int8)
				*val[i] = int8(bitstream.ReadInt(8))
			}
			params[i] = val
		case 85:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint8, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint8)
				*val[i] = uint8(bitstream.ReadInt(8))
			}
			params[i] = val
		case 86:
			nLen := bitstream.ReadInt(16)
			val := make([]*int16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int16)
				*val[i] = int16(bitstream.ReadInt(16))
			}
			params[i] = val
		case 87:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint16, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint16)
				*val[i] = uint16(bitstream.ReadInt(16))
			}
			params[i] = val
		case 88:
			nLen := bitstream.ReadInt(16)
			val := make([]*int32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int32)
				*val[i] = int32(bitstream.ReadInt(32))
			}
			params[i] = val
		case 89:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint32, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint32)
				*val[i] = uint32(bitstream.ReadInt(32))
			}
			params[i] = val
		case 90:
			nLen := bitstream.ReadInt(16)
			val := make([]*int64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int64)
				*val[i] = int64(bitstream.ReadInt64())
			}
			params[i] = val
		case 91:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint64, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint64)
				*val[i] = uint64(bitstream.ReadInt64())
			}
			params[i] = val
		case 92:
			nLen := bitstream.ReadInt(16)
			val := make([]*string, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(string)
				*val[i] = bitstream.ReadString()
			}
			params[i] = val
		case 93:
			nLen := bitstream.ReadInt(16)
			val := make([]*int, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(int)
				*val[i] = bitstream.ReadInt(32)
			}
			params[i] = val
		case 94:
			nLen := bitstream.ReadInt(16)
			val := make([]*uint, nLen)
			for i := 0; i < nLen; i++ {
				val[i] = new(uint)
				*val[i] = uint(bitstream.ReadInt(32))
			}
			params[i] = val

		case 101:
			nLen := bitstream.ReadInt(16)
			aa := bool(false)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**bool)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEBOOL
				val1 := bitstream.ReadFlag()
				*value = &val1
			}
			params[i] = val.Interface()
		case 102:
			nLen := bitstream.ReadInt(16)
			aa := float64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**float64)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEFLOAT64
				val1 := bitstream.ReadFloat64()
				*value = &val1
			}
			params[i] = val.Interface()
		case 103:
			nLen := bitstream.ReadInt(16)
			aa := float32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**float32)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEFLOAT32
				val1 := float32(bitstream.ReadFloat64())
				*value = &val1
			}
			params[i] = val.Interface()
		case 104:
			nLen := bitstream.ReadInt(16)
			aa := int8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**int8)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEINT8
				val1 := int8(bitstream.ReadInt(8))
				*value = &val1
			}
			params[i] = val.Interface()
		case 105:
			nLen := bitstream.ReadInt(16)
			aa := uint8(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**uint8)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEUINT8
				val1 := uint8(bitstream.ReadInt(8))
				*value = &val1
			}
			params[i] = val.Interface()
		case 106:
			nLen := bitstream.ReadInt(16)
			aa := int16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**int16)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEINT16
				val1 := int16(bitstream.ReadInt(16))
				*value = &val1
			}
			params[i] = val.Interface()
		case 107:
			nLen := bitstream.ReadInt(16)
			aa := uint16(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**uint16)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEUINT16
				val1 := uint16(bitstream.ReadInt(16))
				*value = &val1
			}
			params[i] = val.Interface()
		case 108:
			nLen := bitstream.ReadInt(16)
			aa := int32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**int32)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEINT32
				val1 := int32(bitstream.ReadInt(32))
				*value = &val1
			}
			params[i] = val.Interface()
		case 109:
			nLen := bitstream.ReadInt(16)
			aa := uint32(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**uint32)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEUINT32
				val1 := uint32(bitstream.ReadInt(32))
				*value = &val1
			}
			params[i] = val.Interface()
		case 110:
			nLen := bitstream.ReadInt(16)
			aa := int64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**int64)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEINT64
				val1 := int64(bitstream.ReadInt64())
				*value = &val1
			}
			params[i] = val.Interface()
		case 111:
			nLen := bitstream.ReadInt(16)
			aa := uint64(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**uint64)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEUINT64
				val1 := uint64(bitstream.ReadInt64())
				*value = &val1
			}
			params[i] = val.Interface()
		case 112:
			nLen := bitstream.ReadInt(16)
			aa := string("")
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**string)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZESTRING
				val1 := string(bitstream.ReadString())
				*value = &val1
			}
			params[i] = val.Interface()
		case 113:
			nLen := bitstream.ReadInt(16)
			aa := int(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**int)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEINT
				val1 := bitstream.ReadInt(32)
				*value = &val1
			}
			params[i] = val.Interface()
		case 114:
			nLen := bitstream.ReadInt(16)
			aa := uint(0)
			tVal := reflect.ArrayOf(nLen, reflect.TypeOf(&aa))
			val := reflect.New(tVal).Elem()
			arrayPtr := uintptr(unsafe.Pointer(val.Addr().Pointer()))
			for i := 0; i < nLen; i++ {
				value := (**uint)(unsafe.Pointer(arrayPtr))
				arrayPtr = arrayPtr + SIZEUINT
				val1 := uint(bitstream.ReadInt(32))
				*value = &val1
			}
			params[i] = val.Interface()

		case RPCMESSAGE: //protobuf
			// packet := message.GetPakcetByName(funcName)
			// nLen := bitstream.ReadInt(base.Bit32)
			// packetBuf := bitstream.ReadBits(nLen << 3)
			// message.UnmarshalText(packet, packetBuf)
			// params[i] = packet

		case RPCGOB: //gob
			nLen := bitstream.ReadInt(base.Bit32)
			packetBuf := bitstream.ReadBits(nLen << 3)
			val := reflect.New(funcType.In(i))
			buf := bytes.NewBuffer(packetBuf)
			decoder := gob.NewDecoder(buf)
			err := decoder.DecodeValue(val)
			if err != nil {
				log.Printf("func [%s] params no fit, params [%v], error[%s]", funcType.String(), params, err)
				return nil
			}
			params[i] = val.Elem().Interface()

		default:
			panic("func [%s] params type not supported")
		}
	}
	// }
	return params
	// if k.NumIn() != len(params) {
	// 	log.Printf("func [%s] can not call, func params [%s], params [%v]", funcName, strParams, params)
	// 	return
	// }
}
