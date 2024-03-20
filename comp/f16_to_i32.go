package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func F16_to_i32(a Float16_t, exact bool) int32 {
	var uiA uint16
	var sign bool
	var exp int8
	var frac uint16
	var sig32 int32
	var shiftDist int8

	uiA = uint16(a)
	sign = SignF16UI(uiA)
	exp = ExpF16UI(uiA)
	frac = FracF16UI(uiA)

	if exp == 0x1F {
		if frac != 0 {
			return 0x7FFFFFFF
		} else {
			if sign {
				return -0x7FFFFFFF - 1
			} else {
				return 0x7FFFFFFF
			}
		}
	}

	sig32 = int32(frac)
	if exp != 0 {
		sig32 |= 0x0400
		shiftDist = exp - 0x19
		if shiftDist >= 0 {
			sig32 <<= shiftDist
			if sign {
				return -sig32
			} else {
				return sig32
			}
		}
		shiftDist = exp - 0x0D
		if shiftDist > 0 {
			sig32 <<= shiftDist
		}
	}
	return Softfloat_roundToI32(sign, uint64(uint32(sig32)), exact)
}
