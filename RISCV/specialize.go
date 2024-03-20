package riscv

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_isSigNaNF16UI(uiA uint16) bool {
	return (((uiA & 0x7E00) == 0x7C00) && IntToBool(int(uiA&0x01FF)))
}
