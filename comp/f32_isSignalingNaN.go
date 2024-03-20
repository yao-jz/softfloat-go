package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
	. "github.com/yao-jz/softfloat-go/riscv"
)

func F32_isSignalingNaN(a Float32_t) bool {
	uiA := uint32(a)
	return Softfloat_isSigNaNF32UI(uiA)
}
