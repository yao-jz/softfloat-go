package comp

func Softfloat_roundToI64(sign bool, sig uint64, sigExtra uint64, exact bool) int64 {
	var ui uint64
	var z int64
	if uint64(0x8000000000000000) < sigExtra {
		sig += 1
		if sig == 0 {
			if sign {
				return (-int64(0x7FFFFFFFFFFFFFFF) - 1)
			} else {
				return int64(0x7FFFFFFFFFFFFFFF)
			}
		}
		if sigExtra == uint64(0x8000000000000000) {
			sig &= ^uint64(1)
		}
	}
	if sign {
		ui = uint64(-int64(sig))
	} else {
		ui = sig
	}
	z = int64(ui)
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
			return (-int64(0x7FFFFFFFFFFFFFFF) - 1)
		} else {
			return int64(0x7FFFFFFFFFFFFFFF)
		}
	}
	return z
}
