package delta

import "errors"

type reader struct {
	a   []rune // source array
	pos int    // current position in array
}

func haveRune(reader *reader) bool {
	return reader.pos < len(reader.a)
}

func getRune(reader *reader) (rune, error) {
	b := reader.a[reader.pos]
	reader.pos++
	if reader.pos > len(reader.a) {
		return -1, errors.New("out of bounds")
	}
	return b, nil
}
func getChar(reader *reader) (string, error) {
	value, error := getRune(reader)
	return string(value), error
}

// Read base64-encoded unsigned integer.
func getInt(reader *reader) (int, error) {
	v := 0
	for haveRune(reader) {
		value, error := getRune(reader)
		if error != nil {
			return -1, error
		}
		c := zValue[0x7f&value]
		if c < 0 {
			break
		}
		v = (v << 6) + c
	}
	reader.pos--
	return v >> 0, nil
}
