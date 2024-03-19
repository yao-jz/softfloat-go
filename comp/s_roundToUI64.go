package comp

func Softfloat_roundToUI64(sign bool, sig uint64, sigExtra uint64, exact bool) uint64 {
	if uint64(0x8000000000000000) <= sigExtra {
		// increment
		sig += 1
		if sig == 0 {
			// invalid
			if sign {
				return 0
			} else {
				return uint64(0xFFFFFFFFFFFFFFFF)
			}
		}
		if sigExtra == uint64(0x8000000000000000) {
			sig &= ^uint64(1)
		}
	}
	if sign && (sig != 0) {
		// invalid
		if sign {
			return 0
		} else {
			return uint64(0xFFFFFFFFFFFFFFFF)
		}
	}
	return sig
}
