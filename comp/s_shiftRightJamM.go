package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_shiftRightJamM(size_words uint8, aPtr []uint32, dist uint32, zPtr []uint32) {
	var wordJam, wordDist uint32
	var ptr []uint32
	var i uint8
	var innerDist uint8

	wordJam = 0
	wordDist = dist >> 5
	if wordDist != 0 {
		if uint32(size_words) < wordDist {
			wordDist = uint32(size_words)
		}
		ptr := aPtr[IndexMultiwordLo(int(size_words), int(wordDist)):]
		i = uint8(wordDist)
		wordJam = ptr[0]
		ptr = ptr[1:]
		if wordJam != 0 {
			// skip the loop
		} else {
			i--
			for i > 0 {
				wordJam = ptr[0]
				ptr = ptr[1:]
				if wordJam != 0 {
					break
				}
				i--
			}
		}
		ptr = zPtr
	}
	if wordDist < uint32(size_words) {
		aPtr = aPtr[uint(IndexMultiwordHiBut(int(size_words), int(wordDist))):]
		innerDist = uint8(dist & 31)
		if innerDist != 0 {
			Softfloat_shortShiftRightJamM(
				uint8(uint32(size_words)-wordDist),
				aPtr,
				uint8(innerDist),
				zPtr[IndexMultiwordLoBut(int(size_words), int(wordDist)):],
			)
			if wordDist == 0 {
				// wordJam
				if wordJam != 0 {
					zPtr[IndexWordLo(int(size_words))] |= 1
				}
			}
		} else {
			aPtr = aPtr[IndexWordLo(int(uint32(size_words)-wordDist)):]
			ptr := zPtr[IndexWordLo(int(size_words)):]
			for i = uint8(uint32(size_words) - wordDist); i > 0; i-- {
				ptr[0] = aPtr[0]
				aPtr = aPtr[1:]
				ptr = ptr[1:]
			}
		}
		ptr = zPtr[IndexMultiwordHi(int(size_words), int(wordDist)):]
	}
	ptr[0] = 0
	ptr = ptr[1:]
	wordDist--
	for wordDist > 0 {
		ptr[0] = 0
		ptr = ptr[1:]
		wordDist--
	}
	// wordJam:
	//
	if wordJam != 0 {
		zPtr[IndexWordLo(int(size_words))] |= 1
	}
}
