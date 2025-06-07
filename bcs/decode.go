package bcs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Unmarshal(data []byte, v any) (int, error) {
	return NewDecoder(bytes.NewReader(data)).Decode(v)
}

type Decoder struct {
	reader     io.Reader
	byteBuffer [1]byte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (d *Decoder) Decode(v any) (int, error) {
	reflectValue := reflect.ValueOf(v)
	if reflectValue.Kind() != reflect.Pointer || reflectValue.IsNil() {
		return 0, fmt.Errorf("not a pointer or nil pointer")
	}

	return d.decode(reflectValue)
}

func (d *Decoder) decode(v reflect.Value) (int, error) {
	if !v.CanInterface() {
		return 0, nil
	}

	if i, isUnmarshaler := v.Interface().(Unmarshaler); isUnmarshaler {
		return i.UnmarshalBCS(d.reader)
	}

	if _, isEnum := v.Interface().(Enum); isEnum {
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			return d.decodeEnum(v.Elem())
		default:
			return d.decodeEnum(v)
		}
	}

	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return d.decode(v.Elem())

	case reflect.Interface:
		if v.IsNil() {
			return 0, fmt.Errorf("cannot decode into nil interface")
		}
		return d.decode(v.Elem())

	case reflect.Chan, reflect.Func, reflect.Uintptr, reflect.UnsafePointer:

		return 0, nil
	default:
		return d.decodeVanilla(v)
	}
}
func (d *Decoder) decodeVanilla(v reflect.Value) (int, error) {
	kind := v.Kind()
	if !v.CanSet() {
		return 0, fmt.Errorf("cannot change value of kind %s", kind.String())
	}

	switch kind {
	case reflect.Bool:
		t, n, err := d.readByte()
		if err != nil {
			return n, err
		}

		if t == 0 {
			v.SetBool(false)
		} else {
			v.SetBool(true)
		}

		return n, nil

	case reflect.Int8, reflect.Uint8:
		return 1, binary.Read(d.reader, binary.LittleEndian, v.Addr().Interface())
	case reflect.Int16, reflect.Uint16:
		return 2, binary.Read(d.reader, binary.LittleEndian, v.Addr().Interface())
	case reflect.Int32, reflect.Uint32:
		return 4, binary.Read(d.reader, binary.LittleEndian, v.Addr().Interface())
	case reflect.Int64, reflect.Uint64:
		return 8, binary.Read(d.reader, binary.LittleEndian, v.Addr().Interface())

	case reflect.Struct:
		return d.decodeStruct(v)

	case reflect.Slice:
		sliceType := v.Type().Elem()
		if sliceType.Kind() == reflect.Uint8 {
			return d.decodeByteSlice(v)
		}

		return d.decodeSlice(v)

	case reflect.Array:
		arrayType := v.Type().Elem()
		if arrayType.Kind() == reflect.Uint8 {
			return d.decodeByteArray(v)
		}
		return d.decodeArray(v)

	case reflect.String:
		return d.decodeString(v)

	default:
		return 0, fmt.Errorf("unsupported vanilla decoding type: %s", kind.String())
	}
}

func (d *Decoder) decodeString(v reflect.Value) (int, error) {
	size, n, err := ULEB128Decode[int](d.reader)
	if err != nil {
		return n, err
	}

	if size == 0 {
		v.SetString("")
		return n, nil
	}

	tmp := make([]byte, size)

	read, err := d.reader.Read(tmp)
	n += read
	if err != nil {
		return n, err
	}

	if size != read {
		return n, fmt.Errorf("wrong number of bytes read for string, want: %d, got %d", size, read)
	}

	v.SetString(string(tmp))

	return n, nil
}

func (d *Decoder) readByte() (byte, int, error) {
	b := d.byteBuffer[:]
	n, err := d.reader.Read(b)
	if err != nil {
		return 0, n, err
	}
	if n == 0 {
		return 0, n, io.ErrUnexpectedEOF
	}

	return b[0], n, nil
}

func (d *Decoder) decodeStruct(v reflect.Value) (int, error) {
	t := v.Type()

	var n int

fieldLoop:
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanInterface() {
			continue fieldLoop
		}
		tag, err := parseTagValue(t.Field(i).Tag.Get(tagName))
		if err != nil {
			return n, err
		}

		switch {
		case tag.isIgnored():
			continue fieldLoop
		case tag.isOptional():
			isOptional, k, err := d.readByte()
			n += k
			if err != nil {
				return n, err
			}
			if isOptional == 0 {
				field.Set(reflect.Zero(field.Type()))
			} else {
				field.Set(reflect.New(field.Type().Elem()))
				k, err := d.decode(field.Elem())
				n += k
				if err != nil {
					return n, err
				}
			}
		default:
			k, err := d.decode(field)
			n += k
			if err != nil {
				return n, err
			}
		}
	}

	return n, nil
}

func (d *Decoder) decodeEnum(v reflect.Value) (int, error) {
	if v.Kind() != reflect.Struct {
		return 0, fmt.Errorf("only support struct for Enum, got %s", v.Kind().String())
	}
	enumId, n, err := ULEB128Decode[int](d.reader)
	if err != nil {
		return n, err
	}

	field := v.Field(enumId)

	k, err := d.decode(field)
	n += k

	return n, err
}

func (d *Decoder) decodeByteSlice(v reflect.Value) (int, error) {
	size, n, err := ULEB128Decode[int](d.reader)
	if err != nil {
		return n, err
	}

	if size == 0 {
		return n, nil
	}

	tmp := make([]byte, size)

	read, err := d.reader.Read(tmp)
	n += read
	if err != nil {
		return n, err
	}

	if size != read {
		return n, fmt.Errorf("wrong number of bytes read for []byte, want: %d, got %d", size, read)
	}

	v.Set(reflect.ValueOf(tmp))

	return n, nil
}

func (d *Decoder) decodeByteArray(v reflect.Value) (int, error) {
	arraySize := v.Len()

	if arraySize == 0 {
		return 0, nil
	}

	tmp := make([]byte, arraySize)

	read, err := d.reader.Read(tmp)
	if err != nil {
		return read, err
	}

	if arraySize != read {
		return read, fmt.Errorf("wrong number of bytes read for [%d]byte, want: %d, got %d", arraySize, arraySize, read)
	}

	for i := 0; i < arraySize; i++ {
		v.Index(i).SetUint(uint64(tmp[i]))
	}

	return read, nil
}

func (d *Decoder) decodeArray(v reflect.Value) (int, error) {
	size := v.Len()
	t := v.Type()
	elementType := t.Elem()

	var n int
	if elementType.Kind() == reflect.Pointer {
		for i := 0; i < size; i++ {
			idx := reflect.New(elementType.Elem())
			k, err := d.decode(idx.Elem())
			n += k
			if err != nil {
				return n, err
			}
			v.Index(i).Set(idx)
		}
	} else {
		for i := 0; i < size; i++ {
			idx := reflect.New(elementType)
			k, err := d.decode(idx.Elem())
			n += k
			if err != nil {
				return n, err
			}
			v.Index(i).Set(idx.Elem())
		}
	}

	return n, nil
}

func (d *Decoder) decodeSlice(v reflect.Value) (int, error) {

	size, n, err := ULEB128Decode[int](d.reader)
	if err != nil {
		return n, err
	}

	elementType := v.Type().Elem()

	tmp := reflect.MakeSlice(v.Type(), 0, size)

	if elementType.Kind() == reflect.Pointer {
		for i := 0; i < size; i++ {
			ind := reflect.New(elementType.Elem())
			k, err := d.decode(ind)
			n += k
			if err != nil {
				return n, err
			}
			tmp = reflect.Append(tmp, ind)
		}
	} else {
		for i := 0; i < size; i++ {
			ind := reflect.New(elementType)
			k, err := d.decode(ind.Elem())
			n += k
			if err != nil {
				return n, err
			}
			tmp = reflect.Append(tmp, ind.Elem())
		}
	}

	v.Set(tmp)

	return n, nil
}
