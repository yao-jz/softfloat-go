package comp

func Softfloat_roundToI32(sign bool, sig uint64, exact bool) int32 {
	var roundIncrement, roundBits uint16
	var sig32 uint32
	var z int32
	var ui uint32
	roundIncrement = 0x800

	roundBits = uint16(sig & 0xFFF)
	sig += uint64(roundIncrement)
	if sig&uint64(0xFFFFF00000000000) != 0 {
		if sign {
			return (-0x7FFFFFFF - 1)
		} else {
			return 0x7FFFFFFF
		}
	}
	sig32 = uint32(sig >> 12)
	if roundBits == 0x800 {
		sig32 &= ^uint32(1)
	}
	if sign {
		ui = uint32(-int(sig32))
	} else {
		ui = sig32
	}
	z = int32(ui)
	zlz := int32(0)
	if z < 0 {
		zlz = int32(1)
	}
	signt := int32(0)
	if sign {
		signt = int32(1)
	}
	zt := false
	if z != 0 {
		zt = true
	}
	xor := false
	if (zlz ^ signt) != 0 {
		xor = true
	}
	if zt && xor {
		if sign {
			return (-0x7FFFFFFF - 1)
		} else {
			return 0x7FFFFFFF
		}
	}
	return z
}
