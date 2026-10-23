package comp_test

import (
	"math"
	"testing"

	"github.com/yao-jz/softfloat-go/comp"
	. "github.com/yao-jz/softfloat-go/include"
)

// FCLASS returns a one-hot mask: bit 0 -inf, 1 negative normal, 2 negative subnormal, 3 -0,
// 4 +0, 5 positive subnormal, 6 positive normal, 7 +inf, 8 signaling NaN, 9 quiet NaN.
func TestClassify(t *testing.T) {
	cases64 := []struct {
		bits uint64
		bit  uint
	}{
		{math.Float64bits(math.Inf(-1)), 0},
		{math.Float64bits(-1.5), 1},
		{math.Float64bits(-math.SmallestNonzeroFloat64), 2},
		{1 << 63, 3},
		{0, 4},
		{0x000FFFFFFFFFFFFF, 5}, // largest subnormal
		{math.Float64bits(0x1p-1022), 6},
		{math.Float64bits(math.MaxFloat64), 6},
		{math.Float64bits(math.Inf(1)), 7},
		{sNaN64, 8},
		{0xFFF0000000000001, 8}, // negative signaling NaN
		{canonicalNaN64, 9},
		{qNaN64, 9},
	}
	for _, c := range cases64 {
		if got := comp.F64_classS(Float64_t(c.bits)); got != 1<<c.bit {
			t.Errorf("F64_classS(%#016x) = %#x, want %#x", c.bits, got, 1<<c.bit)
		}
	}
	cases32 := []struct {
		bits uint32
		bit  uint
	}{
		{math.Float32bits(float32(math.Inf(-1))), 0},
		{math.Float32bits(-1.5), 1},
		{0x80000001, 2},
		{1 << 31, 3},
		{0, 4},
		{0x007FFFFF, 5},
		{math.Float32bits(0x1p-126), 6},
		{math.Float32bits(float32(math.Inf(1))), 7},
		{sNaN32, 8},
		{canonicalNaN32, 9},
		{qNaN32, 9},
	}
	for _, c := range cases32 {
		if got := comp.F32_classS(Float32_t(c.bits)); got != 1<<c.bit {
			t.Errorf("F32_classS(%#08x) = %#x, want %#x", c.bits, got, 1<<c.bit)
		}
	}
}

// FMIN and FMAX (RISC-V 2.2 and later) return the other operand when exactly one input is
// NaN, the canonical NaN when both are, treat -0 as less than +0, and raise invalid only for a
// signaling NaN.
func TestMinMax(t *testing.T) {
	negZero, posZero := uint64(1<<63), uint64(0)
	one, two := math.Float64bits(1), math.Float64bits(2)
	cases := []struct {
		name     string
		a, b     uint64
		min, max uint64
		flags    uint8
	}{
		{"1, 2", one, two, one, two, 0},
		{"2, 1", two, one, one, two, 0},
		{"-0, +0", negZero, posZero, negZero, posZero, 0},
		{"+0, -0", posZero, negZero, negZero, posZero, 0},
		{"qNaN, 1", qNaN64, one, one, one, 0},
		{"1, qNaN", one, qNaN64, one, one, 0},
		{"sNaN, 1", sNaN64, one, one, one, nv},
		{"qNaN, qNaN", qNaN64, qNaN64, canonicalNaN64, canonicalNaN64, 0},
		{"sNaN, qNaN", sNaN64, qNaN64, canonicalNaN64, canonicalNaN64, nv},
		{"-inf, 1", math.Float64bits(math.Inf(-1)), one, math.Float64bits(math.Inf(-1)), one, 0},
	}
	for _, c := range cases {
		var min, max Float64_t
		fmin := raised(func() { min = comp.F64_min(Float64_t(c.a), Float64_t(c.b)) })
		fmax := raised(func() { max = comp.F64_max(Float64_t(c.a), Float64_t(c.b)) })
		if uint64(min) != c.min || fmin != c.flags {
			t.Errorf("F64_min(%s) = %#x flags %05b, want %#x flags %05b", c.name, uint64(min), fmin, c.min, c.flags)
		}
		if uint64(max) != c.max || fmax != c.flags {
			t.Errorf("F64_max(%s) = %#x flags %05b, want %#x flags %05b", c.name, uint64(max), fmax, c.max, c.flags)
		}
	}

	neg0, pos0 := uint32(1<<31), uint32(0)
	one32 := math.Float32bits(1)
	cases32 := []struct {
		name     string
		a, b     uint32
		min, max uint32
		flags    uint8
	}{
		{"-0, +0", neg0, pos0, neg0, pos0, 0},
		{"qNaN, 1", qNaN32, one32, one32, one32, 0},
		{"sNaN, 1", sNaN32, one32, one32, one32, nv},
		{"qNaN, qNaN", qNaN32, qNaN32, canonicalNaN32, canonicalNaN32, 0},
	}
	for _, c := range cases32 {
		var min, max Float32_t
		fmin := raised(func() { min = comp.F32_min(Float32_t(c.a), Float32_t(c.b)) })
		fmax := raised(func() { max = comp.F32_max(Float32_t(c.a), Float32_t(c.b)) })
		if uint32(min) != c.min || fmin != c.flags {
			t.Errorf("F32_min(%s) = %#x flags %05b, want %#x flags %05b", c.name, uint32(min), fmin, c.min, c.flags)
		}
		if uint32(max) != c.max || fmax != c.flags {
			t.Errorf("F32_max(%s) = %#x flags %05b, want %#x flags %05b", c.name, uint32(max), fmax, c.max, c.flags)
		}
	}
}

// Out-of-range and NaN conversions saturate and raise invalid, as RISC-V's FCVT requires:
// NaN and +overflow give the largest value, -overflow the smallest (0 for unsigned).
func TestConversionEdges(t *testing.T) {
	inf, nan := f64(math.Inf(1)), f64(math.NaN())
	cases := []struct {
		name  string
		got   func() uint64
		want  uint64
		flags uint8
	}{
		{"i32(NaN)", func() uint64 { return uint64(uint32(comp.F64_to_i32(nan, true))) }, 0x7FFFFFFF, nv},
		{"i32(+inf)", func() uint64 { return uint64(uint32(comp.F64_to_i32(inf, true))) }, 0x7FFFFFFF, nv},
		{"i32(-3e9)", func() uint64 { return uint64(uint32(comp.F64_to_i32(f64(-3e9), true))) }, 0x80000000, nv},
		{"u32(-1)", func() uint64 { return uint64(comp.F64_to_ui32(f64(-1), true)) }, 0, nv},
		{"u32(NaN)", func() uint64 { return uint64(comp.F64_to_ui32(nan, true)) }, 0xFFFFFFFF, nv},
		{"i64(NaN)", func() uint64 { return uint64(comp.F64_to_i64(nan, true)) }, 0x7FFFFFFFFFFFFFFF, nv},
		{"i64(-inf)", func() uint64 { return uint64(comp.F64_to_i64(f64(math.Inf(-1)), true)) }, 1 << 63, nv},
		{"i64(2^63)", func() uint64 { return uint64(comp.F64_to_i64(f64(0x1p63), true)) }, 0x7FFFFFFFFFFFFFFF, nv},
		{"u64(2^64)", func() uint64 { return comp.F64_to_ui64(f64(0x1p64), true) }, 0xFFFFFFFFFFFFFFFF, nv},
		{"u64(-1)", func() uint64 { return comp.F64_to_ui64(f64(-1), true) }, 0, nv},
		{"u64(-0.4)", func() uint64 { return comp.F64_to_ui64(f64(-0.4), true) }, 0, nx}, // rounds to -0: in range
		{"u64(NaN)", func() uint64 { return comp.F64_to_ui64(nan, true) }, 0xFFFFFFFFFFFFFFFF, nv},
		{"i64(2.5), exact", func() uint64 { return uint64(comp.F64_to_i64(f64(2.5), true)) }, 2, nx},
		{"i64(2.5), not exact", func() uint64 { return uint64(comp.F64_to_i64(f64(2.5), false)) }, 2, 0},
		{"f32 i32(NaN)", func() uint64 { return uint64(uint32(comp.F32_to_i32(f32(float32(math.NaN())), true))) }, 0x7FFFFFFF, nv},
		{"f32 u64(-1)", func() uint64 { return comp.F32_to_ui64(f32(-1), true) }, 0, nv},
	}
	for _, c := range cases {
		var got uint64
		flags := raised(func() { got = c.got() })
		if got != c.want || flags != c.flags {
			t.Errorf("%s: got %#x flags %05b, want %#x flags %05b", c.name, got, flags, c.want, c.flags)
		}
	}
}

// Conversions to integers honour the rounding mode.
func TestConversionRounding(t *testing.T) {
	want := map[uint8][2]int64{ // results for 2.5 and -2.5
		Softfloat_round_near_even:   {2, -2},
		Softfloat_round_minMag:      {2, -2},
		Softfloat_round_min:         {2, -3},
		Softfloat_round_max:         {3, -2},
		Softfloat_round_near_maxMag: {3, -3},
	}
	for mode, w := range want {
		setRounding(t, mode)
		for i, v := range []float64{2.5, -2.5} {
			if got := comp.F64_to_i64(f64(v), false); got != w[i] {
				t.Errorf("F64_to_i64(%v) in mode %d = %d, want %d", v, mode, got, w[i])
			}
			if got := comp.F32_to_i32(f32(float32(v)), false); int64(got) != w[i] {
				t.Errorf("F32_to_i32(%v) in mode %d = %d, want %d", v, mode, got, w[i])
			}
		}
	}
}

// Integral values in range must convert exactly in both directions. This is the check that
// catches a conversion ignoring its operand (F64_to_ui64 once returned 0 for every input).
func TestConversionsExact(t *testing.T) {
	integral := []float64{0, 1, -1, 2, -2, 3, 100, -100, 12345, -12345, 1 << 20, 1 << 30,
		1 << 52, -(1 << 52), 2147483647, -2147483648, 4294967295, 1e15, -1e15, 1e18}
	for _, v := range integral {
		bits := f64(v)
		if got := comp.F64_to_i64(bits, true); got != int64(v) {
			t.Errorf("F64_to_i64(%v) = %d", v, got)
		}
		if v >= 0 {
			if got := comp.F64_to_ui64(bits, true); got != uint64(v) {
				t.Errorf("F64_to_ui64(%v) = %d", v, got)
			}
		}
		if v >= math.MinInt32 && v <= math.MaxInt32 {
			if got := comp.F64_to_i32(bits, true); got != int32(v) {
				t.Errorf("F64_to_i32(%v) = %d", v, got)
			}
		}
		if v >= 0 && v <= math.MaxUint32 {
			if got := comp.F64_to_ui32(bits, true); got != uint32(v) {
				t.Errorf("F64_to_ui32(%v) = %d", v, got)
			}
		}
		if f := float32(v); float64(f) == v && v >= math.MinInt32 && v <= math.MaxInt32 {
			if got := comp.F32_to_i32(f32(f), true); got != int32(v) {
				t.Errorf("F32_to_i32(%v) = %d", v, got)
			}
		}
	}
	for _, v := range []int64{0, 1, -1, 2, -2, 12345, -12345, 1 << 31, -(1 << 31), 1 << 52} {
		if got := uint64(comp.I64_to_f64(v)); got != math.Float64bits(float64(v)) {
			t.Errorf("I64_to_f64(%d) = %#x", v, got)
		}
		if got := uint64(comp.UI64ToF64(uint64(v))); got != math.Float64bits(float64(uint64(v))) {
			t.Errorf("UI64ToF64(%d) = %#x", uint64(v), got)
		}
		if v >= math.MinInt32 && v <= math.MaxInt32 {
			if got := uint64(comp.I32_to_f64(int32(v))); got != math.Float64bits(float64(v)) {
				t.Errorf("I32_to_f64(%d) = %#x", v, got)
			}
			if got := uint32(comp.I32_to_f32(int32(v))); got != math.Float32bits(float32(v)) {
				t.Errorf("I32_to_f32(%d) = %#x", v, got)
			}
		}
	}
	// And float-to-float: every float32 widens exactly, and narrowing rounds to nearest.
	for _, v := range []float32{0, 1, -1.5, math.MaxFloat32, math.SmallestNonzeroFloat32, 0x1p-126} {
		if got := uint64(comp.F32_to_f64(f32(v))); got != math.Float64bits(float64(v)) {
			t.Errorf("F32_to_f64(%v) = %#x", v, got)
		}
	}
	for _, v := range []float64{0.1, -1e-40, 3.4028235e38, 1e39, 1.0000000596046448} {
		if got := uint32(comp.F64_to_f32(f64(v))); got != math.Float32bits(float32(v)) {
			t.Errorf("F64_to_f32(%v) = %#x, want %#x", v, got, math.Float32bits(float32(v)))
		}
	}
}
