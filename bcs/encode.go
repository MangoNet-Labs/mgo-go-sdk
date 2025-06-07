package bcs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

type Marshaler interface {
	MarshalBCS() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalBCS(io.Reader) (int, error)
}

type Enum interface {
	IsBcsEnum()
}

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (e *Encoder) Encode(v any) error {
	return e.encode(reflect.ValueOf(v))
}

func (e *Encoder) encode(v reflect.Value) error {

	if !v.IsValid() || !v.CanInterface() {
		return nil
	}

	i := v.Interface()
	if m, ismarshaler := i.(Marshaler); ismarshaler {
		bytes, err := m.MarshalBCS()
		if err != nil {
			return err
		}

		_, err = e.w.Write(bytes)

		return err
	}
	if _, isenum := i.(Enum); isenum {
		return e.encodeEnum(reflect.Indirect(v))
	}

	kind := v.Kind()

	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return binary.Write(e.w, binary.LittleEndian, v.Interface())

	case reflect.Pointer:

		return e.encode(reflect.Indirect(v))

	case reflect.Interface:
		return e.encode(v.Elem())

	case reflect.Slice:
		if byteSlice, ok := (v.Interface()).([]byte); ok {
			return e.encodeByteSlice(byteSlice)
		}
		return e.encodeSlice(v)

	case reflect.Array:
		if byteSlice := fixedByteArrayToSlice(v); byteSlice != nil {
			return e.encodeByteSlice(byteSlice)
		}
		return e.encodeArray(v)

	case reflect.String:
		str := []byte(v.String())
		return e.encodeByteSlice(str)

	case reflect.Struct:
		return e.encodeStruct(v)

	case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.UnsafePointer:
		return nil

	default:
		return fmt.Errorf("unsupported kind: %s, consider make the field ignored by using - tag or provide a customized Marshaler implementation", kind.String())
	}
}

func (e *Encoder) encodeEnum(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}

		fieldType := t.Field(i)
		// check the tag
		tag, err := parseTagValue(fieldType.Tag.Get(tagName))
		if err != nil {
			return err
		}
		if tag.isIgnored() {
			continue
		}
		fieldKind := field.Kind()
		if fieldKind != reflect.Pointer && fieldKind != reflect.Interface {
			return fmt.Errorf("enum only supports fields that are either pointers or interfaces, unless they are ignored")
		}
		if !field.IsNil() {
			if _, err := e.w.Write(ULEB128Encode(i)); err != nil {
				return err
			}
			if fieldKind == reflect.Pointer {
				return e.encode(reflect.Indirect(field))
			} else {
				return e.encode(v)
			}
		}
	}

	return fmt.Errorf("no field is set in the enum")
}

func (e *Encoder) encodeByteSlice(b []byte) error {
	l := len(b)
	if _, err := e.w.Write(ULEB128Encode(l)); err != nil {
		return err
	}

	if _, err := e.w.Write(b); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeArray(v reflect.Value) error {
	length := v.Len()
	for i := 0; i < length; i++ {
		if err := e.encode(v.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeSlice(v reflect.Value) error {
	length := v.Len()
	if _, err := e.w.Write(ULEB128Encode(length)); err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		if err := e.encode(v.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeStruct(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue
		}
		tag, err := parseTagValue(t.Field(i).Tag.Get(tagName))
		if err != nil {
			return err
		}
		switch {
		case tag.isIgnored():
			continue
		case tag.isOptional():
			if field.Kind() != reflect.Pointer && field.Kind() != reflect.Interface {
				return fmt.Errorf("optional field can only be pointer or interface")
			}
			if field.IsNil() {
				_, err := e.w.Write([]byte{0})
				if err != nil {
					return err
				}
			} else {
				if _, err := e.w.Write([]byte{1}); err != nil {
					return err
				}
				if err := e.encode(field.Elem()); err != nil {
					return err
				}
			}
			continue
		default:
			if err := e.encode(field); err != nil {
				return err
			}
		}
	}

	return nil
}

func Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	e := NewEncoder(&b)

	if err := e.Encode(v); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

type Option[T any] struct {
	Some T
	None bool
}

func (p *Option[T]) MarshalBCS() ([]byte, error) {
	if p.None {
		return []byte{0}, nil
	}
	b, err := Marshal(p.Some)
	return append([]byte{1}, b...), err
}

func (p *Option[T]) UnmarshalBCS(r io.Reader) (int, error) {
	buf := new(bytes.Buffer)
	io.Copy(buf, r)
	tmp := buf.Bytes()
	if len(tmp) == 1 {
		p.None = true
		return 1, nil
	}
	b := tmp[1:]
	return Unmarshal(b, &p.Some)
}

func MustMarshal(v any) []byte {
	result, err := Marshal(v)
	if err != nil {
		panic(err)
	}

	return result
}

func fixedByteArrayToSlice(v reflect.Value) []byte {
	if v.Kind() != reflect.Array {
		return nil
	}
	if v.Type().Elem().Kind() != reflect.Uint8 {
		return nil
	}

	if isCustomMgoAddressBytes(v) {
		return nil
	}

	length := v.Len()
	slice := make([]byte, length)
	for i := 0; i < length; i++ {
		slice[i] = byte(v.Index(i).Uint())
	}
	return slice
}

const MgoAddressBytesName = "MgoAddressBytes"

func isCustomMgoAddressBytes(v reflect.Value) bool {
	return v.Type().Name() == MgoAddressBytesName
}
