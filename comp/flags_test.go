package comp_test

import (
	"math"
	"testing"

	"github.com/yao-jz/softfloat-go/comp"
	. "github.com/yao-jz/softfloat-go/include"
)

const (
	nx = Softfloat_flag_inexact
	uf = Softfloat_flag_underflow
	of = Softfloat_flag_overflow
	dz = Softfloat_flag_infinite // RISC-V calls this divide by zero
	nv = Softfloat_flag_invalid

	sNaN64 = 0x7FF0000000000001
	qNaN64 = 0x7FF8000000000001 // a quiet NaN with a payload
	sNaN32 = 0x7F800001
	qNaN32 = 0x7FC00001
)

func f64(v float64) Float64_t { return Float64_t(math.Float64bits(v)) }
func f32(v float32) Float32_t { return Float32_t(math.Float32bits(v)) }

// raised runs fn with the flags cleared and returns the flags it raised.
func raised(fn func()) uint8 {
	Softfloat_exceptionFlags = 0
	fn()
	return Softfloat_exceptionFlags
}

func TestExceptionFlags64(t *testing.T) {
	inf, max := math.Inf(1), math.MaxFloat64
	cases := []struct {
		name  string
		op    func() Float64_t
		want  uint64
		flags uint8
	}{
		{"exact", func() Float64_t { return comp.F64_add(f64(1), f64(2)) }, math.Float64bits(3), 0},
		{"inexact", func() Float64_t { return comp.F64_add(f64(1), f64(0x1p-60)) }, math.Float64bits(1), nx},
		{"0 × ∞", func() Float64_t { return comp.F64_mul(f64(0), f64(inf)) }, canonicalNaN64, nv},
		{"∞ - ∞", func() Float64_t { return comp.F64_sub(f64(inf), f64(inf)) }, canonicalNaN64, nv},
		{"0 / 0", func() Float64_t { return comp.F64_div(f64(0), f64(0)) }, canonicalNaN64, nv},
		{"1 / 0", func() Float64_t { return comp.F64_div(f64(1), f64(0)) }, math.Float64bits(inf), dz},
		{"-1 / 0", func() Float64_t { return comp.F64_div(f64(-1), f64(0)) }, math.Float64bits(-inf), dz},
		{"sqrt(-1)", func() Float64_t { return comp.F64_sqrt(f64(-1)) }, canonicalNaN64, nv},
		{"sqrt(-0)", func() Float64_t { return comp.F64_sqrt(f64(math.Copysign(0, -1))) }, 1 << 63, 0},
		{"overflow", func() Float64_t { return comp.F64_mul(f64(max), f64(2)) }, math.Float64bits(inf), of | nx},
		{"tiny, exact", func() Float64_t { return comp.F64_mul(f64(0x1p-1022), f64(0.5)) }, 1 << 51, 0},
		{"tiny, inexact", func() Float64_t { return comp.F64_mul(f64(0x1.0000000000001p-1022), f64(0x1p-3)) },
			math.Float64bits(0x1p-1025), uf | nx},
		{"0 × ∞ + qNaN", func() Float64_t { return comp.F64_mulAdd(f64(0), f64(inf), qNaN64) }, canonicalNaN64, nv},
		{"quiet NaN operand", func() Float64_t { return comp.F64_add(qNaN64, f64(1)) }, canonicalNaN64, 0},
		{"signaling NaN operand", func() Float64_t { return comp.F64_add(sNaN64, f64(1)) }, canonicalNaN64, nv},
		{"signaling NaN, mul", func() Float64_t { return comp.F64_mul(f64(1), sNaN64) }, canonicalNaN64, nv},
	}
	for _, c := range cases {
		var got Float64_t
		flags := raised(func() { got = c.op() })
		if uint64(got) != c.want || flags != c.flags {
			t.Errorf("%s: got %#x flags %05b, want %#x flags %05b", c.name, uint64(got), flags, c.want, c.flags)
		}
	}
}

func TestExceptionFlags32(t *testing.T) {
	inf := float32(math.Inf(1))
	cases := []struct {
		name  string
		op    func() Float32_t
		want  uint32
		flags uint8
	}{
		{"inexact", func() Float32_t { return comp.F32_div(f32(1), f32(3)) }, math.Float32bits(float32(1) / 3), nx},
		{"0 × ∞", func() Float32_t { return comp.F32_mul(f32(0), f32(inf)) }, canonicalNaN32, nv},
		{"1 / 0", func() Float32_t { return comp.F32_div(f32(1), f32(0)) }, math.Float32bits(inf), dz},
		{"sqrt(-1)", func() Float32_t { return comp.F32_sqrt(f32(-1)) }, canonicalNaN32, nv},
		{"overflow", func() Float32_t { return comp.F32_mul(f32(math.MaxFloat32), f32(2)) }, math.Float32bits(inf), of | nx},
		{"tiny, inexact", func() Float32_t { return comp.F32_mul(f32(0x1.000002p-126), f32(0x1p-3)) },
			math.Float32bits(0x1p-129), uf | nx},
		{"signaling NaN operand", func() Float32_t { return comp.F32_sub(sNaN32, f32(1)) }, canonicalNaN32, nv},
		{"quiet NaN operand", func() Float32_t { return comp.F32_sub(qNaN32, f32(1)) }, canonicalNaN32, 0},
	}
	for _, c := range cases {
		var got Float32_t
		flags := raised(func() { got = c.op() })
		if uint32(got) != c.want || flags != c.flags {
			t.Errorf("%s: got %#x flags %05b, want %#x flags %05b", c.name, uint32(got), flags, c.want, c.flags)
		}
	}
}

// FEQ is a quiet comparison: only a signaling NaN is invalid. FLT and FLE are signaling: any
// NaN is invalid. Every comparison with a NaN is false.
func TestComparisons(t *testing.T) {
	type cmp struct {
		name  string
		got   func() bool
		want  bool
		flags uint8
	}
	one, two := f64(1), f64(2)
	negZero, posZero := f64(math.Copysign(0, -1)), f64(0)
	cases := []cmp{
		{"1 == 1", func() bool { return comp.F64_eq(one, one) }, true, 0},
		{"-0 == +0", func() bool { return comp.F64_eq(negZero, posZero) }, true, 0},
		{"-0 < +0", func() bool { return comp.F64_lt(negZero, posZero) }, false, 0},
		{"-0 <= +0", func() bool { return comp.F64_le(negZero, posZero) }, true, 0},
		{"1 < 2", func() bool { return comp.F64_lt(one, two) }, true, 0},
		{"2 <= 1", func() bool { return comp.F64_le(two, one) }, false, 0},
		{"qNaN == 1", func() bool { return comp.F64_eq(qNaN64, one) }, false, 0},
		{"sNaN == 1", func() bool { return comp.F64_eq(sNaN64, one) }, false, nv},
		{"qNaN < 1", func() bool { return comp.F64_lt(qNaN64, one) }, false, nv},
		{"1 <= qNaN", func() bool { return comp.F64_le(one, qNaN64) }, false, nv},
		{"f32 qNaN == 1", func() bool { return comp.F32_eq(qNaN32, f32(1)) }, false, 0},
		{"f32 sNaN == 1", func() bool { return comp.F32_eq(sNaN32, f32(1)) }, false, nv},
		{"f32 qNaN < 1", func() bool { return comp.F32_lt(qNaN32, f32(1)) }, false, nv},
		{"f32 1 <= qNaN", func() bool { return comp.F32_le(f32(1), qNaN32) }, false, nv},
		{"f32 -1 < 1", func() bool { return comp.F32_lt(f32(-1), f32(1)) }, true, 0},
	}
	for _, c := range cases {
		var got bool
		flags := raised(func() { got = c.got() })
		if got != c.want || flags != c.flags {
			t.Errorf("%s: got %v flags %05b, want %v flags %05b", c.name, got, flags, c.want, c.flags)
		}
	}
}
