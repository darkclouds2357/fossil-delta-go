package delta

import (
	"errors"
	"math"
)

// digitCount : Return the number digits in the base64 representation of a positive integer.
func digitCount(v int) int {
	var i, x int = 1, 64
	for ; v >= x; i, x = i+1, x<<6 {
		/* nothing */
	}
	return i
}

// Return a 32-bit checksum of the array.
func checksum(arr []rune) int {
	var sum0, sum1, sum2, sum3 rune = 0, 0, 0, 0
	z := 0
	N := len(arr)

	for N >= 16 {
		sum0 = sum0 + arr[z+0]
		sum1 = sum1 + arr[z+1]
		sum2 = sum2 + arr[z+2]
		sum3 = sum3 + arr[z+3]

		sum0 = sum0 + arr[z+4]
		sum1 = sum1 + arr[z+5]
		sum2 = sum2 + arr[z+6]
		sum3 = sum3 + arr[z+7]

		sum0 = sum0 + arr[z+8]
		sum1 = sum1 + arr[z+9]
		sum2 = sum2 + arr[z+10]
		sum3 = sum3 + arr[z+11]

		sum0 = sum0 + arr[z+12]
		sum1 = sum1 + arr[z+13]
		sum2 = sum2 + arr[z+14]
		sum3 = sum3 + arr[z+15]

		z += 16
		N -= 16
	}
	for N >= 4 {
		sum0 = sum0 + arr[z+0]
		sum1 = sum1 + arr[z+1]
		sum2 = sum2 + arr[z+2]
		sum3 = sum3 + arr[z+3]
		z += 4
		N -= 4
	}
	sum3 += (sum2 << 8) + (sum1 << 16) + (sum0 << 24)
	/* jshint32 -W086 */
	switch N {
	case 3:
		sum3 += arr[z+2] << 8 /* falls through */
	case 2:
		sum3 += arr[z+1] << 16 /* falls through */
	case 1:
		sum3 += arr[z+0] << 24 /* falls through */
	}
	return int(sum3 >> 0)
}

//
// Create : a new delta.
//
// The delta is written into a preallocated buffer, zDelta, which
// should be at least 60 bytes longer than the target file, zOut.
// The delta string will be NUL-terminated, but it might also contain
// embedded NUL characters if either the zSrc or zOut files are
// binary.  This function returns the length of the delta string
// in bytes, excluding the final NUL terminator character.
//
// Output Format:
//
// The delta begins with a base64 number followed by a newline.  This
// number is the number of bytes in the TARGET file.  Thus, given a
// delta file z, a program can compute the size of the output file
// simply by reading the first line and decoding the base-64 number
// found there.  The delta_output_size() routine does exactly this.
//
// After the initial size number, the delta consists of a series of
// literal text segments and commands to copy from the SOURCE file.
// A copy command looks like this:
//
//     NNN@MMM,
//
// where NNN is the number of bytes to be copied and MMM is the offset
// into the source file of the first byte (both base-64).   If NNN is 0
// it means copy the rest of the input file.  Literal text is like this:
//
//     NNN:TTTTT
//
// where NNN is the number of bytes of text (base-64) and TTTTT is the text.
//
// The last term is of the form
//
//     NNN;
//
// In this case, NNN is a 32-bit bigendian checksum of the output file
// that can be used to verify that the delta applied correctly.  All
// numbers are in base-64.
//
// Pure text files generate a pure text delta.  Binary files generate a
// delta that may contain some binary data.
//
// Algorithm:
//
// The encoder first builds a hash table to help it find matching
// patterns in the source file.  16-byte chunks of the source file
// sampled at evenly spaced intervals are used to populate the hash
// table.
//
// Next we begin scanning the target file using a sliding 16-byte
// window.  The hash of the 16-byte window in the target is used to
// search for a matching section in the source file.  When a match
// is found, a copy command is added to the delta.  An effort is
// made to extend the matching section to regions that come before
// and after the 16-byte hash window.  A copy command is only issued
// if the result would use less space that just quoting the text
// literally. Literal text is added to the delta for sections that
// do not match or which can not be encoded efficiently using copy
// commands.
//
func Create(orgin string, target string) []rune {
	var zDelta writer

	lenTarget := len(target)
	lenOrgin := len(orgin)
	lastRead := -1
	// var i
	putInt(&zDelta, lenTarget)
	putChar(&zDelta, "\n")

	// If the source is very small, it means that we have no
	// chance of ever doing a copy command.  Just output a single
	// literal segment for the entire target and exit.
	if lenOrgin <= NHASH {
		putInt(&zDelta, lenTarget)
		putChar(&zDelta, ":")
		putArray(&zDelta, []rune(target), 0, lenTarget)
		putInt(&zDelta, checksum([]rune(target)))
		putChar(&zDelta, ";")
		return toArray(&zDelta)
	}

	// Compute the hash table used to locate matching sections in the source.
	nHash := int(math.Ceil(float64(lenOrgin) / float64(NHASH))) /* Number of hash table entries */
	collide := make([]int, nHash)
	landmark := make([]int, nHash)
	for i := 0; i < nHash; i++ {
		collide[i] = -1
		landmark[i] = -1
	}
	var hv int
	var h rollingHash
	i := 0
	for i = 0; i < lenOrgin-NHASH; i += NHASH {
		initRollingHash(&h, []rune(orgin), i)
		hv = value(&h) % nHash
		collide[i/NHASH] = landmark[hv]
		landmark[hv] = i / NHASH
	}
	base := 0
	var iSrc, iBlock, bestCnt, bestOfst, bestLitsz int

	for base+NHASH < lenTarget {
		bestOfst = 0
		bestLitsz = 0
		initRollingHash(&h, []rune(target), base)
		i = 0 // Trying to match a landmark against zOut[base+i]
		bestCnt = 0
		for {
			limit := 250
			hv = value(&h) % nHash
			iBlock = landmark[hv]
			for iBlock >= 0 && limit > 0 {
				limit--
				//
				// The hash window has identified a potential match against
				// landmark block iBlock.  But we need to investigate further.
				//
				// Look for a region in zOut that matches zSrc. Anchor the search
				// at zSrc[iSrc] and zOut[base+i].  Do not include anything prior to
				// zOut[base] or after zOut[outLen] nor anything after zSrc[srcLen].
				//
				// Set cnt equal to the length of the match and set ofst so that
				// zSrc[ofst] is the first element of the match.  litsz is the number
				// of characters between zOut[base] and the beginning of the match.
				// sz will be the overhead (in bytes) needed to encode the copy
				// command.  Only generate copy command if the overhead of the
				// copy command is less than the amount of literal text to be copied.
				//
				var cnt, ofst, litsz int
				var j, k, x, y int
				var sz int

				// Beginning at iSrc, match forwards as far as we can.
				// j counts the number of characters that match.
				iSrc = iBlock * NHASH

				for j, x, y = 0, iSrc, base+i; x < lenOrgin && y < lenTarget; j, x, y = j+1, x+1, y+1 {
					if orgin[x] != target[y] {
						break
					}
				}
				j--

				// Beginning at iSrc-1, match backwards as far as we can.
				// k counts the number of characters that match.
				for k = 1; k < iSrc && k <= i; k++ {
					if orgin[iSrc-k] != target[base+i-k] {
						break
					}
				}
				k--

				// Compute the offset and size of the matching region.
				ofst = iSrc - k
				cnt = j + k + 1
				litsz = i - k // Number of bytes of literal text before the copy
				// sz will hold the number of bytes needed to encode the "insert"
				// command and the copy command, not counting the "insert" text.
				sz = digitCount(i-k) + digitCount(cnt) + digitCount(ofst) + 3
				if cnt >= sz && cnt > bestCnt {
					// Remember this match only if it is the best so far and it
					// does not increase the file size.
					bestCnt = cnt
					bestOfst = iSrc - k
					bestLitsz = litsz
				}

				// Check the next matching block
				iBlock = collide[iBlock]
			}

			// We have a copy command that does not cause the delta to be larger
			// than a literal insert.  So add the copy command to the delta.
			if bestCnt > 0 {
				if bestLitsz > 0 {
					// Add an insert command before the copy.
					putInt(&zDelta, bestLitsz)
					putChar(&zDelta, ":")
					putArray(&zDelta, []rune(target), base, base+bestLitsz)
					base += bestLitsz
				}
				base += bestCnt
				putInt(&zDelta, bestCnt)
				putChar(&zDelta, "@")
				putInt(&zDelta, bestOfst)
				putChar(&zDelta, ",")
				if bestOfst+bestCnt-1 > lastRead {
					lastRead = bestOfst + bestCnt - 1
				}
				bestCnt = 0
				break
			}

			// If we reach this point, it means no match is found so far
			if base+i+NHASH >= lenTarget {
				// We have reached the end and have not found any
				// matches.  Do an "insert" for everything that does not match
				putInt(&zDelta, lenTarget-base)
				putChar(&zDelta, ":")
				putArray(&zDelta, []rune(target), base, base+lenTarget-base)
				base = lenTarget
				break
			}
			// Advance the hash by one character. Keep looking for a match.
			next(&h, int(target[base+i+NHASH]))
			i++
		}
	}
	// Output a final "insert" record to get all the text at the end of
	// the file that does not match anything in the source.
	if base < lenTarget {
		putInt(&zDelta, lenTarget-base)
		putChar(&zDelta, ":")
		putArray(&zDelta, []rune(target), base, base+lenTarget-base)
	}
	// Output the final checksum record.
	v := checksum([]rune(target))
	putInt(&zDelta, v)
	putChar(&zDelta, ";")
	return toArray(&zDelta)
}

// OutputSize :
// Return the size (in bytes) of the output from applying
// a delta.
//
// This routine is provided so that an procedure that is able
// to call delta_apply() can learn how much space is required
// for the output and hence allocate nor more space that is really
// needed.
//
func OutputSize(delta []rune) (int, error) {
	zDelta := reader{delta, 0}
	size, intError := getInt(&zDelta)
	if intError != nil {
		return -1, intError
	}
	deltaChar, charError := getChar(&zDelta)
	if charError != nil {
		return -1, charError
	}
	if deltaChar != "\n" {
		return -1, errors.New("size integer not terminated by \"\\n\"")
	}
	return size, nil
}

//
// Apply : a delta.
//
// The output buffer should be big enough to hold the whole output
// file and a NUL terminator at the end.  The delta_output_size()
// routine will determine this size for you.
//
// The delta string should be null-terminated.  But the delta string
// may contain embedded NUL characters (if the input and output are
// binary files) so we also have to pass in the length of the delta in
// the lenDelta parameter.
//
// This function returns the size of the output file in bytes (excluding
// the final NUL terminator character).  Except, if the delta string is
// malformed or intended for use with a source file other than zSrc,
// then this routine returns -1.
//
// Refer to the delta_create() documentation above for a description
// of the delta file format.
//
func Apply(orgin string, delta []rune, verifyChecksum bool) ([]rune, error) {
	total := 0
	zDelta := reader{delta, 0}
	lenSrc := len(orgin)
	lenDelta := len(delta)

	limit, intError1 := getInt(&zDelta)
	if intError1 != nil {
		return []rune{-1}, intError1
	}
	char, charError := getChar(&zDelta)
	if charError != nil {
		return []rune{-1}, charError
	}
	if char != "\n" {
		return []rune{-1}, errors.New("size integer not terminated by \"\\n\"")
	}
	zOut := writer{}
	for haveRune(&zDelta) {
		var ofst int
		cnt, intError2 := getInt(&zDelta)
		if intError2 != nil {
			return []rune{-1}, intError2
		}
		loopChar, loopCharError := getChar(&zDelta)
		if loopCharError != nil {
			return []rune{-1}, loopCharError
		}
		switch loopChar {
		case "@":
			ofst, intError2 = getInt(&zDelta)
			if intError2 != nil {
				return []rune{-1}, intError2
			}
			getDeltaChar, caseCharError := getChar(&zDelta)
			if caseCharError != nil {
				return []rune{-1}, caseCharError
			}
			if haveRune(&zDelta) && getDeltaChar != "," {
				return []rune{-1}, errors.New("copy command not terminated by \",\"")
			}
			total += cnt
			if total > limit {
				return []rune{-1}, errors.New("copy exceeds output file size")
			}
			if ofst+cnt > lenSrc {
				return []rune{-1}, errors.New("copy extends past end of input")
			}
			putArray(&zOut, []rune(orgin), ofst, ofst+cnt)
			break

		case ":":
			total += cnt
			if total > limit {
				return []rune{-1}, errors.New("insert command gives an output larger than predicted")
			}
			if cnt > lenDelta {
				return []rune{-1}, errors.New("insert count exceeds size of delta")
			}
			putArray(&zOut, zDelta.a, zDelta.pos, zDelta.pos+cnt)
			zDelta.pos += cnt
			break

		case ";":
			out := toArray(&zOut)
			if verifyChecksum && cnt != checksum(out) {
				return []rune{-1}, errors.New("bad checksum")
			}
			if total != limit {
				return []rune{-1}, errors.New("generated size does not match predicted size")
			}
			return out, nil

		default:
			return []rune{-1}, errors.New("unknown delta operator")
		}
	}
	return []rune{-1}, errors.New("unterminated delta")
}
