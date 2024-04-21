package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

type BinarySerializer struct {
	output []byte
}

func NewBinarySerializer() *BinarySerializer {
	return &BinarySerializer{output: make([]byte, 0)}
}

type Serializable interface {
	Call(interface{})
}

type Uint8List []byte

type Int128 struct {
	Low  int64
	High int64
}

type Uint128 struct {
	Low  uint64
	High uint64
}

// * Bool
func (bs *BinarySerializer) SerializeBool(val bool) {
	if val {
		bs.output = append(bs.output, 1)
	} else {
		bs.output = append(bs.output, 0)
	}
}

// * Option tag
func (bs *BinarySerializer) SerializeOptionTag(val bool) {
	if val {
		bs.output = append(bs.output, 1)
	} else {
		bs.output = append(bs.output, 0)
	}
}

// * Uint8
func (bs *BinarySerializer) SerializeUint8(val uint8) {
	bs.output = append(bs.output, val)
}

// * Uint16
func (bs *BinarySerializer) SerializeUint16(val uint16) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Uint32
func (bs *BinarySerializer) SerializeUint32(val uint32) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Uint64
func (bs *BinarySerializer) SerializeUint64(val uint64) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Uint64 another solution
func (bs *BinarySerializer) SerializeUint64_2(val uint64) {
	buf := make([]byte, 8)
	for i := 0; i < 8; i++ {
		buf[i] = byte(val >> (i * 8) & 0xFF)
	}
	bs.output = append(bs.output, buf...)
}

// * Int8
func (bs *BinarySerializer) SerializeInt8(val int8) {
	bs.output = append(bs.output, byte(val))
}

// * Int16
func (bs *BinarySerializer) SerializeInt16(val int16) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Int32
func (bs *BinarySerializer) SerializeInt32(val int32) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Int64
func (bs *BinarySerializer) SerializeInt64(val int64) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Float32
func (bs *BinarySerializer) SerializeFloat32(val float32) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Float64
func (bs *BinarySerializer) SerializeFloat64(val float64) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, val)
	bs.output = append(bs.output, buf.Bytes()...)
}

// * Char
func (bs *BinarySerializer) SerializeChar(val int32) {
	bs.SerializeInt32(val)
}

// * Uint8List
func (bs *BinarySerializer) SerializeUint8List(val Uint8List) {
	bs.SerializeLengthUint30(len(val))
	bs.output = append(bs.output, val...)
}

// * Bytes
func (bs *BinarySerializer) SerializeBytes(val []byte) {
	bs.SerializeLengthUint30(len(val))
	bs.output = append(bs.output, val...)
}

// * String
func (bs *BinarySerializer) SerializeString(str string) {
	bs.SerializeUint8List(Uint8List(str))
}

// * Uint30
func (bs *BinarySerializer) SerializeLengthUint30(number int) {
	b4 := byte(number & 0xFF)
	if number < (1 << 6) {
		bs.output = append(bs.output, b4)
		return
	}
	b3 := byte((number >> 8) & 0xFF)
	if number < (1 << 14) {
		bs.output = append(bs.output, 0x40|b3, b4)
		return
	}
	b2 := byte((number >> 16) & 0xFF)
	if number < (1 << 22) {
		bs.output = append(bs.output, 0x80|b2, b3, b4)
		return
	}
	b1 := byte((number >> 24) & 0xFF)
	if number < (1 << 30) {
		bs.output = append(bs.output, 0xC0|b1, b2, b3, b4)
		return
	}
	panic("out of range integral type conversion attempted")
}

// * Uint15
func (bs *BinarySerializer) SerializeLengthUint15(number int) {
	b2 := byte(number)
	if number < (1 << 7) {
		bs.SerializeUint8(b2)
		return
	}
	if number < (1 << 15) {
		b1 := byte((number >> 8) & 0xFF)
		bs.output = append(bs.output, 0x80|b1, b2)
		return
	}
	panic("out of range integral type conversion attempted")
}

// * Int128
func (bs *BinarySerializer) SerializeInt128(value Int128) {
	bs.SerializeInt64(value.Low)
	bs.SerializeInt64(value.High)
}

// * Uint128
func (bs *BinarySerializer) SerializeUint128(value Uint128) {
	bs.SerializeUint64(value.Low)
	bs.SerializeUint64(value.High)
}

// * Leb128
func (bs *BinarySerializer) SerializeLeb128(value int, isSigned bool) {
	parts := make([]byte, 0)
	if isSigned {
		for {
			byteVal := value & 0x7F
			value >>= 7
			if (value == 0 && byteVal&0x40 == 0) || (value == -1 && byteVal&0x40 > 0) {
				parts = append(parts, byte(byteVal))
				break
			} else {
				parts = append(parts, byte(byteVal|0x80))
			}
		}
	} else {
		size := int(math.Ceil(float64(len(strconv.FormatInt(int64(value), 2)) / 7.0)))
		for i := 0; i < size; i++ {
			part := value & 0x7F
			value >>= 7
			parts = append(parts, byte(part))
		}
		for i := 0; i < len(parts)-1; i++ {
			parts[i] |= 0x80
		}
	}
	bs.output = append(bs.output, parts...)
}

// * Result
// func (bs *BinarySerializer) SerializeResult(ok Serializable, err Serializable) func(interface{}) error {
// 	return func(obj interface{}) error {
// 		switch obj.(type) {
// 		case Ok:
// 			bs.SerializeUint8(1)
// 			ok.Call(obj)
// 		case Err:
// 			bs.SerializeUint8(0)
// 			err.Call(obj)
// 		default:
// 			return errors.New("Wrong type, probably not supported")
// 		}
// 		return nil
// 	}
// }

// * Nullable
func (bs *BinarySerializer) SerializeNullable(input Serializable) func(interface{}) error {
	return func(object interface{}) error {
		// Serializing option tag
		bs.SerializeOptionTag(object != nil)
		if object != nil {
			input.Call(object)
		}
		return nil
	}
}

// * Fixed List
func (bs *BinarySerializer) SerializeFixedList(input Serializable) func([]interface{}) {
	return func(objects []interface{}) {
		for _, obj := range objects {
			input.Call(obj)
		}
	}
}

// * List
func (bs *BinarySerializer) SerializeList(input Serializable) func([]interface{}) {
	return func(objects []interface{}) {
		bs.SerializeLengthUint30(len(objects))
		bs.SerializeFixedList(input)(objects)
	}
}

func main() {
	bs := NewBinarySerializer()
	bs.SerializeInt8(127)
	fmt.Print(bs.output)
}
