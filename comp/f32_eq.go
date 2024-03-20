package comp

import (
	. "github.com/yao-jz/softfloat-go/include"
)

func F32_eq(a, b Float32_t) bool {
	uiA := uint32(a)
	uiB := uint32(b)

	if IsNaNF32UI(uiA) || IsNaNF32UI(uiB) {
		return false
	}
	return (uiA == uiB) || uint32((uiA|uiB)<<1) == 0
}
