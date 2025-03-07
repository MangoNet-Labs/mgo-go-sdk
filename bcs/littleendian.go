package bcs

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
)

var BCS Bcs

type Bcs struct{}

type Uint128 struct {
	High uint64
	Low  uint64
}
type Uint256 struct {
	HighHigh uint64
	HighLow  uint64
	LowHigh  uint64
	LowLow   uint64
}

func (Bcs) Uint8(b []byte) uint8 {
	_ = b[0]
	return uint8(b[0])
}
func (Bcs) Uint16(b []byte) uint16 {
	_ = b[1]
	return uint16(b[0]) | uint16(b[1])<<8
}
func (Bcs) Uint32(b []byte) uint32 {
	_ = b[3]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}
func (Bcs) Uint64(b []byte) uint64 {
	_ = b[7]
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}
func (Bcs) Uint128(b []byte) Uint128 {
	_ = b[15]
	low := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	high := uint64(b[8]) | uint64(b[9])<<8 | uint64(b[10])<<16 | uint64(b[11])<<24 |
		uint64(b[12])<<32 | uint64(b[13])<<40 | uint64(b[14])<<48 | uint64(b[15])<<56
	return Uint128{High: high, Low: low}
}

func (Bcs) Uint256(b []byte) Uint256 {
	_ = b[31]
	lowLow := uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
	lowHigh := uint64(b[8]) | uint64(b[9])<<8 | uint64(b[10])<<16 | uint64(b[11])<<24 |
		uint64(b[12])<<32 | uint64(b[13])<<40 | uint64(b[14])<<48 | uint64(b[15])<<56
	highLow := uint64(b[16]) | uint64(b[17])<<8 | uint64(b[18])<<16 | uint64(b[19])<<24 |
		uint64(b[20])<<32 | uint64(b[21])<<40 | uint64(b[22])<<48 | uint64(b[23])<<56
	highHigh := uint64(b[24]) | uint64(b[25])<<8 | uint64(b[26])<<16 | uint64(b[27])<<24 |
		uint64(b[28])<<32 | uint64(b[29])<<40 | uint64(b[30])<<48 | uint64(b[31])<<56
	return Uint256{HighHigh: highHigh, HighLow: highLow, LowHigh: lowHigh, LowLow: lowLow}
}

func ULEBEncode(num uint64) []byte {
	var encoded []byte

	for {
		byteVal := uint8(num & 0x7F)
		num >>= 7

		if num != 0 {
			byteVal |= 0x80
		}

		encoded = append(encoded, byteVal)

		if num == 0 {
			break
		}
	}

	return encoded
}

func ULEBDecode(encoded []byte) (uint64, int) {
	var num uint64
	var k int
	for i, byteVal := range encoded {
		k = i
		num |= uint64(byteVal&0x7F) << (7 * uint(i))
		if byteVal&0x80 == 0 {
			break
		}
	}
	return num, k
}

func Bcsde(obj any, data []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("provided object is not a pointer to a struct")
	}
	v = v.Elem()

	var index uint
	index = 0

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i).Type
		fieldName := v.Type().Field(i).Name
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set", fieldName)
		}
		switch strings.ToLower(fieldType.Name()) {
		case "address":
			address := hex.EncodeToString(data[index : index+32])
			field.SetString(address)
			index += 32
		case "uint8":
			value := BCS.Uint8(data[index : index+1])
			field.SetUint(uint64(value))
			index++
		case "uint16":
			value := BCS.Uint16(data[index : index+2])
			field.SetUint(uint64(value))
			index += 2
		case "uint32":
			value := BCS.Uint32(data[index : index+4])
			field.SetUint(uint64(value))
			index += 4
		case "uint64":
			value := BCS.Uint64(data[index : index+8])
			field.SetUint(value)
			index += 8
		case "Uint128":
			value := BCS.Uint128(data[index : index+16])
			field.Set(reflect.ValueOf(value))
			index += 16
		case "Uint256":
			value := BCS.Uint256(data[index : index+32])
			field.Set(reflect.ValueOf(value))
			index += 32
		case "bool":
			value := data[index] != 0
			field.SetBool(value)
			index++
		case "string":
			fieldlen, i := ULEBDecode(data[index:])
			index = index + uint(i)
			valuebyte := data[index : index+uint(fieldlen)]
			value := string(valuebyte)
			field.SetString(value)
			index = index + uint(fieldlen)
		}
	}
	return nil
}
