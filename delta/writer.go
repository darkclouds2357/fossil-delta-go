package delta

type writer struct {
	a []rune
}

func toArray(writer *writer) []rune {
	return writer.a
}

func putRune(writer *writer, b rune) {
	writer.a = append(writer.a, b)
}

func putChar(writer *writer, s string) {
	putRune(writer, rune(s[0]))
}

func putInt(writer *writer, v int) {
	var zBuf []rune
	if v == 0 {
		putChar(writer, "0")
		return
	}
	i, j := 0, 0
	for i = 0; v > 0; i, v = i+1, v>>6 {
		zBuf = append(zBuf, zDigits[v&0x3f])
	}
	for j = i - 1; j >= 0; j-- {
		putRune(writer, zBuf[j])
	}
}

// Copy from array at start to end.
func putArray(writer *writer, a []rune, start, end int) {
	for i := start; i < end; i++ {
		writer.a = append(writer.a, a[i])
	}
}
