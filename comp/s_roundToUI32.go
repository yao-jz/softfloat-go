package comp

func Softfloat_roundToUI32(sign bool, sig uint64, exact bool) uint32 {
	var roundIncrement, roundBits uint16
	var z uint32
	roundIncrement = 0x800

	roundBits = uint16(sig & 0xFFF)
	sig += uint64(roundIncrement)
	if sig&uint64(0xFFFFF00000000000) != 0 {
		// invalid
		if sign {
			return 0
		} else {
			return 0xFFFFFFFF
		}
	}
	z = uint32(sig >> 12)
	if roundBits == 0x800 {
		z &= ^uint32(1)
	}
	if sign && (z != 0) {
		// invalid
		if sign {
			return 0
		} else {
			return 0xFFFFFFFF
		}
	}
	return z
}
