package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
	. "github.com/yao-jz/softfloat-go/riscv"
)

func F32_div(a Float32_t, b Float32_t) Float32_t {
	var uiA, uiB uint32
	var signA, signB, signZ bool
	var expA, expB, expZ int16
	var sigA, sigB, sigZ uint32
	var uiZ uint32
	var rem uint64

	uiA = uint32(a)
	signA = SignF32UI(uiA)
	expA = ExpF32UI(uiA)
	sigA = FracF32UI(uiA)
	uiB = uint32(b)
	signB = SignF32UI(uiB)
	expB = ExpF32UI(uiB)
	sigB = FracF32UI(uiB)
	signZ = signA != signB

	if expA == 0xFF {
		if sigA != 0 {
			// propagateNaN
			uiZ = Softfloat_propagateNaNF32UI(uiA, uiB)
			// uiZ
			return Float32_t(uiZ)
		}
		if expB == 0xFF {
			if sigB != 0 {
				// propagateNaN
				uiZ = Softfloat_propagateNaNF32UI(uiA, uiB)
				// uiZ
				return Float32_t(uiZ)
			}
			// invalid
			uiZ = 0x7FC00000
			// uiZ
			return Float32_t(uiZ)
		}
		// infinity
		uiZ = PackToF32UI(signZ, 0xFF, 0)
		// uiZ
		return Float32_t(uiZ)
	}

	if expB == 0xFF {
		if sigB != 0 {
			// propagateNaN
			uiZ = Softfloat_propagateNaNF32UI(uiA, uiB)
			// uiZ
			return Float32_t(uiZ)
		}
		// zero
		uiZ = PackToF32UI(signZ, 0, 0)
		// uiZ
		return Float32_t(uiZ)
	}

	if expB == 0 {
		if sigB == 0 {
			if uint32(expA)|sigA == 0 {
				// invalid
				uiZ = 0x7FC00000
				// uiZ
				return Float32_t(uiZ)
			}
			// infinity
			uiZ = PackToF32UI(signZ, 0xFF, 0)
			// uiZ
			return Float32_t(uiZ)
		}
		expB, sigB = Softfloat_normSubnormalF32Sig(sigB)
	}
	if expA == 0 {
		if sigA == 0 {
			// zero
			uiZ = PackToF32UI(signZ, 0, 0)
			// uiZ
			return Float32_t(uiZ)
		}
		expA, sigA = Softfloat_normSubnormalF32Sig(sigA)
	}

	expZ = expA - expB + 0x7E
	sigA |= 0x00800000
	sigB |= 0x00800000
	if sigA < sigB {
		expZ--
		sigA <<= 8
	} else {
		sigA <<= 7
	}
	sigB <<= 8
	sigZ = uint32((uint64(sigA) * uint64(Softfloat_approxRecip32_1(sigB))) >> 32)

	sigZ += 2
	if (sigZ & 0x3F) < 2 {
		sigZ &= ^uint32(3)
		rem = (uint64(sigA) << 32) - uint64(sigZ<<1)*uint64(sigB)
		if rem&0x8000000000000000 != 0 {
			sigZ -= 4
		} else {
			if rem != 0 {
				sigZ |= 1
			}
		}
	}
	return Softfloat_roundPackToF32(signZ, expZ, sigZ)
}
