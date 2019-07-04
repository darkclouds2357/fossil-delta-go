package delta

// RollingHash :
// The current state of the rolling hash.
//
// z[] holds the values that have been hashed.  z[] is a circular buffer.
// z[i] is the first entry and z[(i+NHASH-1)%NHASH] is the last entry of
// the window.
//
// Hash.a is the sum of all elements of hash.z[].  Hash.b is a weighted
// sum.  Hash.b is z[i]*NHASH + z[i+1]*(NHASH-1) + ... + z[i+NHASH-1]*1.
// (Each index for z[] should be module NHASH, of course.  The %NHASH operator
// is omitted in the prior expression for brevity.)
//
type rollingHash struct {
	a, b int         /* Hash values */
	i    int         /* Start of the hash window */
	z    [NHASH]rune // the values that have been hashed.
}

func initRollingHash(rolling *rollingHash, z []rune, pos int) {
	var a, b int
	x := 0
	for i := 0; i < NHASH; i++ {
		x = int(z[pos+i])
		a = (a + x) & 0xffff
		b = (b + ((NHASH - i) * x)) & 0xffff
		rolling.z[i] = z[pos+i]
	}
	rolling.a = a & 0xffff
	rolling.b = b & 0xffff
	rolling.i = 0
}

func next(rolling *rollingHash, c int) {
	old := int(rolling.z[rolling.i])
	rolling.z[rolling.i] = rune(c)
	rolling.i = (rolling.i + 1) & (NHASH - 1)
	rolling.a = (rolling.a - old + c) & 0xffff
	rolling.b = (rolling.b - NHASH*old + rolling.a) & 0xffff
}

func value(rolling *rollingHash) int {
	return (rolling.a & 0xffff) | ((rolling.b & 0xffff) << 16)
}

func hasOnce(z []rune) int {
	a := int(z[0])
	b := int(z[0])
	for i := 1; i < NHASH; i++ {
		a += int(z[i])
		b += a
	}
	return a | (b << 16)
}
