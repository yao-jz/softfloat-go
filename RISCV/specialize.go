package riscv

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func Softfloat_isSigNaNF16UI(uiA uint16) bool {
	return (((uiA & 0x7E00) == 0x7C00) && IntToBool(int(uiA&0x01FF)))
}

func Softfloat_isSigNaNF32UI(uiA uint32) bool {
	return (((uiA & 0x7FC00000) == 0x7F800000) && IntToBool(int(uiA&0x003FFFFF)))
}
