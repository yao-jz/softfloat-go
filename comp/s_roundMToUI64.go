package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_roundMToUI64(sign bool, extSigPtr []uint32, exact bool) uint64 {
	var sig uint64
	var sigExtra uint32

	sig = uint64(extSigPtr[IndexWord(3, 2)])<<32 | uint64(extSigPtr[IndexWord(3, 1)])
	sigExtra = extSigPtr[IndexWordLo(3)]
	if sigExtra >= 0x80000000 {
		// increment
		sig++
		if sig == 0 {
			// invalid
			if sign {
				return 0
			} else {
				return 0xFFFFFFFFFFFFFFFF
			}
		}
		if sigExtra == 0x80000000 {
			sig &= ^uint64(1)
		}
		if sign && sig != 0 {
			// invalid
			if sign {
				return 0
			} else {
				return 0xFFFFFFFFFFFFFFFF
			}
		}
		return sig
	}
	if sign && sig != 0 {
		// invalid
		if sign {
			return 0
		} else {
			return 0xFFFFFFFFFFFFFFFF
		}
	}
	return sig
}
