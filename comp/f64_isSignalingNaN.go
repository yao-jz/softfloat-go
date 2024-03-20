package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
	. "github.com/yao-jz/softfloat-go/riscv"
)

func F64_isSignalingNaN(a Float64_t) bool {
	uiA := uint64(a)
	return Softfloat_isSigNaNF64UI(uiA)
}
