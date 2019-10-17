package base

import (
	"encoding/binary"
	"log"
	"math"
	"reflect"
	"strings"
)

func Assert(val bool, desc string) {
	if val == false {
		log.Printf("\nFatal : {%s}", desc)
	}
}

func Int8ToBytes(val int8) []byte {
	buff := make([]byte, 1)
	// binary.LittleEndian.PutUint32(buff, uint32(val))
	//LittleEndian
	buff[0] = byte(val)
	return buff
}

func BytesToInt8(data []byte) int8 {
	// buff := make([]byte, 4)
	// tmp := binary.LittleEndian.Uint16(data)
	// binary.LittleEndian.Uint16()

	return int8(data[0])
}

func IntToBytes(val int) []byte {
	buff := make([]byte, 4)
	// binary.LittleEndian.PutUint32(buff, uint32(val))
	//LittleEndian
	buff[0] = byte(val)
	buff[1] = byte(val >> 8)
	buff[2] = byte(val >> 16)
	buff[3] = byte(val >> 24)

	return buff
}

func Int32ToBytes(val int32) []byte {
	buff := make([]byte, 4)
	// binary.LittleEndian.PutUint32(buff, uint32(val))
	//LittleEndian
	buff[0] = byte(val)
	buff[1] = byte(val >> 8)
	buff[2] = byte(val >> 16)
	buff[3] = byte(val >> 24)

	return buff
}

func BytesToInt(data []byte) int {
	// buff := make([]byte, 4)
	// tmp := binary.LittleEndian.Uint32(data)
	l := len(data)
	tmp := 0
	if l <= 1 {
		tmp = int(BytesToInt8(data))
	} else if l <= 2 {
		tmp = int(binary.LittleEndian.Uint16(data))
	} else {
		tmp = int(binary.LittleEndian.Uint32(data))
	}
	return tmp
}

func BytesToInt32(data []byte) int32 {
	// buff := make([]byte, 4)
	tmp := binary.LittleEndian.Uint32(data)
	return int32(tmp)
}

func Int16ToBytes(val int16) []byte {
	buff := make([]byte, 2)
	binary.LittleEndian.PutUint16(buff, uint16(val))
	return buff
}

func BytesToInt16(data []byte) int16 {
	ret := binary.LittleEndian.Uint16(data)
	return int16(ret)
}

func Int64ToBytes(val int64) []byte {
	buff := make([]byte, 8)
	binary.LittleEndian.PutUint64(buff, uint64(val))
	return buff
}

func BytesToInt64(data []byte) int64 {
	ret := binary.LittleEndian.Uint64(data)
	return int64(ret)
}

func Float32ToBytes(val float32) []byte {
	tmp := math.Float32bits(val)
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, tmp)
	return buff
}

func BytesToFloat32(data []byte) float32 {
	tmp := binary.LittleEndian.Uint32(data)
	return math.Float32frombits(tmp)
}

func Float64ToBytes(val float64) []byte {
	tmp := math.Float64bits(val)
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint64(buff, tmp)
	return buff
}

func BytesToFloat64(data []byte) float64 {
	tmp := binary.LittleEndian.Uint64(data)
	return math.Float64frombits(tmp)
}

func SliceTypeString(sTypeName string) string {
	index := strings.Index(sTypeName, "]")
	if index != -1 {
		sTypeName = sTypeName[index+1:]
	}

	if sTypeName == "bool" || sTypeName == "float64" || sTypeName == "float32" || sTypeName == "int8" ||
		sTypeName == "uint8" || sTypeName == "int16" || sTypeName == "uint16" || sTypeName == "int32" ||
		sTypeName == "uint32" || sTypeName == "int64" || sTypeName == "uint64" || sTypeName == "string" ||
		sTypeName == "int" || sTypeName == "uint" {
		return "[]" + sTypeName
	}

	return "[]struct"

	// return sTypeName
}

func ArrayTypeStringEx(sTypeName string) string {
	index := strings.Index(sTypeName, "]")
	if index != -1 {
		sTypeName = sTypeName[index+1:]
	}

	if sTypeName == "bool" || sTypeName == "float64" || sTypeName == "float32" || sTypeName == "int8" ||
		sTypeName == "uint8" || sTypeName == "int16" || sTypeName == "uint16" || sTypeName == "int32" ||
		sTypeName == "uint32" || sTypeName == "int64" || sTypeName == "uint64" || sTypeName == "string" ||
		sTypeName == "int" || sTypeName == "uint" {
		return "[*]" + sTypeName
	}

	return "[*]struct"
}

func TypeString(param interface{}) string {
	pType := reflect.TypeOf(param)
	sType := ""
	if pType.Kind() == reflect.Ptr {
		sType = "*" + pType.Elem().Kind().String()
		pType = pType.Elem()
	} else if pType.Kind() == reflect.Slice {
		sType = SliceTypeString(pType.String())
	} else if pType.Kind() == reflect.Array {
		sType = SliceTypeString(pType.String())
	} else {
		sType = pType.Kind().String()
	}

	if pType.Kind() == reflect.Struct || pType.Kind() == reflect.Map ||
		sType == "[]*struct" || sType == "[]struct" || sType == "[*]*struct" ||
		sType == "[*]struct" {
		sType = "*gob"
	}

	return sType
}
