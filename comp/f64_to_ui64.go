package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func F64_to_ui64(a Float64_t, exact bool) uint64 {
	var uiA uint64
	var sign bool
	var exp int16
	var sig uint64
	var shiftDist int16

	var extSig [3]uint32

	sign = SignF64UI(uiA)
	exp = ExpF64UI(uiA)
	sig = FracF64UI(uiA)

	if exp != 0 {
		sig |= 0x0010000000000000
	}
	shiftDist = 0x433 - exp
	extSig[IndexWord(3, 0)] = 0
	if shiftDist <= 0 {
		if shiftDist < -11 {
			// invalid
			if exp == 0x7FF && FracF64UI(uiA) != 0 {
				return 0xFFFFFFFFFFFFFFFF
			} else {
				if sign {
					return 0
				} else {
					return 0xFFFFFFFFFFFFFFFF
				}
			}
		}
		sig <<= -uint64(shiftDist)
		extSig[IndexWord(3, 2)] = uint32(sig >> 32)
		extSig[IndexWord(3, 1)] = uint32(sig)
	} else {
		extSig[IndexWord(3, 2)] = uint32(sig >> 32)
		extSig[IndexWord(3, 1)] = uint32(sig)
		Softfloat_shiftRightJamM(3, extSig[:], uint32(shiftDist), extSig[:])
	}
	return Softfloat_roundMToUI64(sign, extSig[:], exact)
}
