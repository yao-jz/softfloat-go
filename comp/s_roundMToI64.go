package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_roundMToI64(sign bool, extSigPtr []uint32, exact bool) int64 {
	var sig uint64
	var sigExtra uint32
	var uiZ uint64
	var z int64

	sig = uint64(extSigPtr[IndexWord(3, 2)])<<32 | uint64(extSigPtr[IndexWord(3, 1)])
	sigExtra = extSigPtr[IndexWordLo(3)]
	if sigExtra >= 0x80000000 {
		// increment
		sig++
		if sig == 0 {
			// invalid
			if sign {
				return -0x7FFFFFFFFFFFFFFF - 1
			} else {
				return 0x7FFFFFFFFFFFFFFF
			}
		}
		if sigExtra == 0x80000000 {
			sig &= ^uint64(1)
		}
		if sign {
			uiZ = -sig
		} else {
			uiZ = sig
		}
		z = int64(uiZ)
		if z != 0 && (z < 0) != sign {
			// invalid
			if sign {
				return -0x7FFFFFFFFFFFFFFF - 1
			} else {
				return 0x7FFFFFFFFFFFFFFF
			}
		}
		return z
	}
	if sign {
		uiZ = -sig
	} else {
		uiZ = sig
	}
	z = int64(uiZ)
	if z != 0 && (z < 0) != sign {
		// invalid
		if sign {
			return -0x7FFFFFFFFFFFFFFF - 1
		} else {
			return 0x7FFFFFFFFFFFFFFF
		}
	}
	return z
}
