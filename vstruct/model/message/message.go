package message

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

type _ = strings.Builder
type _ = unsafe.Pointer

var _ = math.Float32frombits
var _ = math.Float64frombits
var _ = strconv.FormatInt
var _ = strconv.FormatUint
var _ = strconv.FormatFloat
var _ = fmt.Sprint

type Types uint8

const (
	Types_File   Types = 0
	Types_Text   Types = 1
	Types_Binary Types = 2
	Types_Signal Types = 3
)

func (e Types) String() string {
	switch e {
	case Types_File:
		return "File"
	case Types_Text:
		return "Text"
	case Types_Binary:
		return "Binary"
	case Types_Signal:
		return "Signal"
	}
	return ""
}

func (e Types) Match(
	onFile func(),
	onText func(),
	onBinary func(),
	onSignal func(),
) {
	switch e {
	case Types_File:
		onFile()
	case Types_Text:
		onText()
	case Types_Binary:
		onBinary()
	case Types_Signal:
		onSignal()
	}
}

type Message []byte

func (s Message) Type() Types {
	return Types(s[0])
}

func (s Message) Order() uint16 {
	_ = s[2]
	var __v uint16 = uint16(s[1]) |
		uint16(s[2])<<8
	return uint16(__v)
}

func (s Message) Data() []byte {
	_ = s[10]
	var __off0 uint64 = 11
	var __off1 uint64 = uint64(s[3]) |
		uint64(s[4])<<8 |
		uint64(s[5])<<16 |
		uint64(s[6])<<24 |
		uint64(s[7])<<32 |
		uint64(s[8])<<40 |
		uint64(s[9])<<48 |
		uint64(s[10])<<56
	return []byte(s[__off0:__off1])
}

func (s Message) Vstruct_Validate() bool {
	if len(s) < 11 {
		return false
	}

	var __off0 uint64 = 11
	var __off1 uint64 = uint64(s[3]) |
		uint64(s[4])<<8 |
		uint64(s[5])<<16 |
		uint64(s[6])<<24 |
		uint64(s[7])<<32 |
		uint64(s[8])<<40 |
		uint64(s[9])<<48 |
		uint64(s[10])<<56
	var __off2 uint64 = uint64(len(s))
	return __off0 <= __off1 && __off1 <= __off2
}

func (s Message) String() string {
	if !s.Vstruct_Validate() {
		return "Message (invalid)"
	}
	var __b strings.Builder
	__b.WriteString("Message {")
	__b.WriteString("Type: ")
	__b.WriteString(s.Type().String())
	__b.WriteString(", ")
	__b.WriteString("Order: ")
	__b.WriteString(strconv.FormatUint(uint64(s.Order()), 10))
	__b.WriteString(", ")
	__b.WriteString("Data: ")
	__b.WriteString(fmt.Sprint(s.Data()))
	__b.WriteString("}")
	return __b.String()
}

func Serialize_Message(dst Message, Type Types, Order uint16, Data []byte) Message {
	_ = dst[10]
	dst[0] = byte(Type)
	dst[1] = byte(Order)
	dst[2] = byte(Order >> 8)

	var __index = uint64(11)
	__tmp_2 := uint64(len(Data)) + __index
	dst[3] = byte(__tmp_2)
	dst[4] = byte(__tmp_2 >> 8)
	dst[5] = byte(__tmp_2 >> 16)
	dst[6] = byte(__tmp_2 >> 24)
	dst[7] = byte(__tmp_2 >> 32)
	dst[8] = byte(__tmp_2 >> 40)
	dst[9] = byte(__tmp_2 >> 48)
	dst[10] = byte(__tmp_2 >> 56)
	copy(dst[__index:__tmp_2], Data)
	return dst
}

func New_Message(Type Types, Order uint16, Data []byte) Message {
	var __vstruct__size = 11 + len(Data)
	var __vstruct__buf = make(Message, __vstruct__size)
	__vstruct__buf = Serialize_Message(__vstruct__buf, Type, Order, Data)
	return __vstruct__buf
}
