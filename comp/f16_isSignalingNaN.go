package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
	. "github.com/yao-jz/softfloat-go/riscv"
)

func F16_isSignalingNaN(a Float16_t) bool {
	uiA := uint16(a)
	return Softfloat_isSigNaNF16UI(uiA)
}
