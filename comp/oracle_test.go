package comp_test

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/yao-jz/softfloat-go/comp"
	. "github.com/yao-jz/softfloat-go/include"
)

// These tests compare the port against the host FPU, an independent, correctly rounded
// implementation of the same standard. For +, -, *, /, sqrt and fused multiply-add the
// correctly rounded result is unique, so any disagreement is a bug on one side.
//
// Two real bugs in this port were found this way, after every conformance suite had passed:
// F64_div dropped a uint32 truncation from the C and was one ULP off on about 2% of inputs,
// and F64_to_ui64 never read its operand and returned 0 for everything.

const (
	canonicalNaN32 = 0x7FC00000
	canonicalNaN64 = 0x7FF8000000000000
	oracleSamples  = 200_000
)

// operand64 is skewed toward the inputs that break implementations: subnormals, wide exponent
// gaps, exact ties and special values. Uniform 64-bit patterns would be mostly NaN.
func operand64(r *rand.Rand) float64 {
	switch r.Intn(6) {
	case 0:
		return r.NormFloat64()
	case 1:
		return math.Float64frombits(r.Uint64() & 0x7FEFFFFFFFFFFFFF) // any finite magnitude
	case 2:
		return r.NormFloat64() * math.Pow(2, float64(r.Intn(600)-300))
	case 3:
		return math.Float64frombits(uint64(r.Intn(1 << 20))) // subnormal
	case 4:
		return float64(r.Intn(1<<20)) / float64(int64(1)<<uint(r.Intn(20))) // exact binary fractions
	}
	specials := []float64{0, math.Copysign(0, -1), 1, -1, math.Inf(1), math.Inf(-1), math.NaN(),
		math.MaxFloat64, -math.MaxFloat64, math.SmallestNonzeroFloat64, 0x1p-1022}
	return specials[r.Intn(len(specials))]
}

func operand32(r *rand.Rand) float32 {
	switch r.Intn(6) {
	case 0:
		return float32(r.NormFloat64())
	case 1:
		return math.Float32frombits(r.Uint32() & 0x7F7FFFFF)
	case 2:
		return float32(r.NormFloat64() * math.Pow(2, float64(r.Intn(160)-80)))
	case 3:
		return math.Float32frombits(uint32(r.Intn(1 << 23))) // subnormal
	case 4:
		return float32(r.Intn(1<<12)) / float32(int32(1)<<uint(r.Intn(12)))
	}
	specials := []float32{0, float32(math.Copysign(0, -1)), 1, -1, float32(math.Inf(1)),
		float32(math.Inf(-1)), float32(math.NaN()), math.MaxFloat32, -math.MaxFloat32,
		math.SmallestNonzeroFloat32, 0x1p-126}
	return specials[r.Intn(len(specials))]
}

// fma32 is the correctly rounded a*b+c in single precision. Computing it through float64 would
// round twice, so it is done exactly and rounded once.
func fma32(a, b, c float32) float32 {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) || math.IsNaN(float64(c)) ||
		math.IsInf(float64(a), 0) || math.IsInf(float64(b), 0) || math.IsInf(float64(c), 0) {
		return float32(math.FMA(float64(a), float64(b), float64(c))) // exact in these cases
	}
	x := new(big.Float).SetPrec(1000).SetFloat64(float64(a))
	x.Mul(x, new(big.Float).SetFloat64(float64(b)))
	x.Add(x, new(big.Float).SetFloat64(float64(c)))
	if x.Sign() == 0 {
		// Exact zero: IEEE gives +0 in round-to-nearest unless both addends are -0.
		return float32(math.FMA(float64(a), float64(b), float64(c)))
	}
	f, _ := x.Float32() // rounds to nearest, ties to even, subnormals included
	return f
}

func TestF64AgainstHostFPU(t *testing.T) {
	ops := []struct {
		name   string
		soft   func(a, b, c uint64) uint64
		native func(a, b, c float64) float64
	}{
		{"F64_add", func(a, b, _ uint64) uint64 { return uint64(comp.F64_add(Float64_t(a), Float64_t(b))) },
			func(a, b, _ float64) float64 { return a + b }},
		{"F64_sub", func(a, b, _ uint64) uint64 { return uint64(comp.F64_sub(Float64_t(a), Float64_t(b))) },
			func(a, b, _ float64) float64 { return a - b }},
		{"F64_mul", func(a, b, _ uint64) uint64 { return uint64(comp.F64_mul(Float64_t(a), Float64_t(b))) },
			func(a, b, _ float64) float64 { return a * b }},
		{"F64_div", func(a, b, _ uint64) uint64 { return uint64(comp.F64_div(Float64_t(a), Float64_t(b))) },
			func(a, b, _ float64) float64 { return a / b }},
		{"F64_sqrt", func(a, _, _ uint64) uint64 { return uint64(comp.F64_sqrt(Float64_t(a))) },
			func(a, _, _ float64) float64 { return math.Sqrt(a) }},
		{"F64_mulAdd", func(a, b, c uint64) uint64 {
			return uint64(comp.F64_mulAdd(Float64_t(a), Float64_t(b), Float64_t(c)))
		}, math.FMA},
	}
	for _, op := range ops {
		t.Run(op.name, func(t *testing.T) {
			r := rand.New(rand.NewSource(1)) // fixed seed, so a failure reproduces
			for i := 0; i < oracleSamples; i++ {
				x, y, z := operand64(r), operand64(r), operand64(r)
				want := op.native(x, y, z)
				got := op.soft(math.Float64bits(x), math.Float64bits(y), math.Float64bits(z))
				if math.IsNaN(want) {
					if got != canonicalNaN64 {
						t.Fatalf("%s(%#x, %#x, %#x) = %#x, want the canonical NaN",
							op.name, math.Float64bits(x), math.Float64bits(y), math.Float64bits(z), got)
					}
					continue
				}
				if got != math.Float64bits(want) {
					t.Fatalf("%s(%#x, %#x, %#x)\n  softfloat %#016x (%v)\n  host FPU  %#016x (%v)",
						op.name, math.Float64bits(x), math.Float64bits(y), math.Float64bits(z),
						got, math.Float64frombits(got), math.Float64bits(want), want)
				}
			}
		})
	}
}

func TestF32AgainstHostFPU(t *testing.T) {
	ops := []struct {
		name   string
		soft   func(a, b, c uint32) uint32
		native func(a, b, c float32) float32
	}{
		{"F32_add", func(a, b, _ uint32) uint32 { return uint32(comp.F32_add(Float32_t(a), Float32_t(b))) },
			func(a, b, _ float32) float32 { return a + b }},
		{"F32_sub", func(a, b, _ uint32) uint32 { return uint32(comp.F32_sub(Float32_t(a), Float32_t(b))) },
			func(a, b, _ float32) float32 { return a - b }},
		{"F32_mul", func(a, b, _ uint32) uint32 { return uint32(comp.F32_mul(Float32_t(a), Float32_t(b))) },
			func(a, b, _ float32) float32 { return a * b }},
		{"F32_div", func(a, b, _ uint32) uint32 { return uint32(comp.F32_div(Float32_t(a), Float32_t(b))) },
			func(a, b, _ float32) float32 { return a / b }},
		{"F32_sqrt", func(a, _, _ uint32) uint32 { return uint32(comp.F32_sqrt(Float32_t(a))) },
			// Correct: the double-precision root of a float32 rounds to float32 without error.
			func(a, _, _ float32) float32 { return float32(math.Sqrt(float64(a))) }},
		{"F32_mulAdd", func(a, b, c uint32) uint32 {
			return uint32(comp.F32_mulAdd(Float32_t(a), Float32_t(b), Float32_t(c)))
		}, fma32},
	}
	for _, op := range ops {
		t.Run(op.name, func(t *testing.T) {
			r := rand.New(rand.NewSource(1))
			for i := 0; i < oracleSamples; i++ {
				x, y, z := operand32(r), operand32(r), operand32(r)
				want := op.native(x, y, z)
				got := op.soft(math.Float32bits(x), math.Float32bits(y), math.Float32bits(z))
				if math.IsNaN(float64(want)) {
					if got != canonicalNaN32 {
						t.Fatalf("%s(%#x, %#x, %#x) = %#x, want the canonical NaN",
							op.name, math.Float32bits(x), math.Float32bits(y), math.Float32bits(z), got)
					}
					continue
				}
				if got != math.Float32bits(want) {
					t.Fatalf("%s(%#x, %#x, %#x)\n  softfloat %#08x (%v)\n  host FPU  %#08x (%v)",
						op.name, math.Float32bits(x), math.Float32bits(y), math.Float32bits(z),
						got, math.Float32frombits(got), math.Float32bits(want), want)
				}
			}
		})
	}
}
