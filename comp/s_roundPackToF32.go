package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_roundPackToF32(sign bool, exp int16, sig uint32) Float32_t {
	var roundIncrement uint16 = 0x40
	var roundBits uint16 = uint16(sig) & 0x7F
	var uiZ uint32

	if 0xFD <= uint32(exp) {
		if exp < 0 {
			sig = Softfloat_shiftRightJam32(sig, uint16(-exp))
			exp = 0
			roundBits = uint16(sig) & 0x7F
		} else if (0xFD < exp) || (0x80000000 <= sig+uint32(roundIncrement)) {
			uiZ = PackToF32UI(sign, 0xFF, 0)
			if roundIncrement == 0 {
				uiZ -= 1
			}
			return Float32_t(uiZ)
		}

	}
	sig = (sig + uint32(roundIncrement)) >> 7
	var roundNearEven uint32 = 1 // Assuming roundNearEven is 1
	temp := uint32(0)
	if (uint32(roundBits) ^ 0x40) == 0 {
		temp = 1
	}
	sig &= ^(uint32(temp & roundNearEven))
	if sig == 0 {
		exp = 0
	}
	uiZ = PackToF32UI(sign, uint32(exp), sig)
	return Float32_t(uiZ)
}
