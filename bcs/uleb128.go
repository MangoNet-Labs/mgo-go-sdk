package bcs

import (
	"encoding/binary"
	"fmt"
	"io"
)

const MaxUleb128Length = 10

type ULEB128SupportedTypes interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint | ~int8 | ~int16 | ~int32 | ~int64 | ~int
}

func ULEB128Encode[T ULEB128SupportedTypes](input T) []byte {
	result := make([]byte, 10)
	i := binary.PutUvarint(result, uint64(input))
	return result[:i]
}

func ULEB128Decode[T ULEB128SupportedTypes](r io.Reader) (T, int, error) {
	buf := make([]byte, 1)
	var v, shift T
	var n int
	for n < 10 {
		i, err := r.Read(buf)
		if i == 0 {
			return 0, n, fmt.Errorf("zero read in. possible EOF")
		}
		if err != nil {
			return 0, n, err
		}
		n += i

		d := T(buf[0])
		ld := d & 127
		if (ld<<shift)>>shift != ld {
			return v, n, fmt.Errorf("overflow at index %d: %v", n-1, ld)
		}

		ld <<= shift
		v = ld + v
		if v < ld {
			return v, n, fmt.Errorf("overflow after adding index %d: %v %v", n-1, ld, v)
		}
		if d <= 127 {
			return v, n, nil
		}

		shift += 7
	}

	return 0, n, fmt.Errorf("failed to find most significant bytes after reading %d bytes", n)
}
