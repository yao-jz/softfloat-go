package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_roundPackToF64(sign bool, exp int16, sig uint64) Float64_t {
	var roundIncrement uint16 = 0x200
	var roundBits uint16 = uint16(sig) & 0x3FF
	var uiZ uint64

	if 0x7FD <= uint16(exp) {
		if exp < 0 {
			sig = Softfloat_shiftRightJam64(sig, uint32(-exp))
			exp = 0
			roundBits = uint16(sig) & 0x3FF
		} else if (0x7FD < exp) || (0x8000000000000000 <= sig+uint64(roundIncrement)) {
			uiZ = PackToF64UI(sign, 0x7FF, 0)
			if roundIncrement == 0 {
				uiZ -= 1
			}
			return Float64_t(uiZ)
		}

	}
	sig = (sig + uint64(roundIncrement)) >> 10
	var roundNearEven uint64 = 1 // Assuming roundNearEven is 1
	temp := uint64(0)
	if (uint64(roundBits) ^ 0x200) == 0 {
		temp = 1
	}
	sig &= ^(uint64(temp & roundNearEven))
	if sig == 0 {
		exp = 0
	}
	uiZ = PackToF64UI(sign, uint64(exp), sig)
	return Float64_t(uiZ)
}
