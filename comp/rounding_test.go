package comp_test

import (
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/yao-jz/softfloat-go/comp"
	. "github.com/yao-jz/softfloat-go/include"
)

// The host FPU only rounds to nearest, ties to even, so the other four modes are checked
// against math/big: compute the exact result (or one with ~1000 bits, which never sits on a
// rounding boundary of a 24- or 53-bit result unless it is exact), then round once in the
// mode under test. math/big has no subnormals, so only results in the normal range count.

var roundingModes = []struct {
	name string
	soft uint8
	big  big.RoundingMode
}{
	{"near_even", Softfloat_round_near_even, big.ToNearestEven},
	{"minMag", Softfloat_round_minMag, big.ToZero},
	{"min", Softfloat_round_min, big.ToNegativeInf},
	{"max", Softfloat_round_max, big.ToPositiveInf},
	{"near_maxMag", Softfloat_round_near_maxMag, big.ToNearestAway},
}

func setRounding(t *testing.T, mode uint8) {
	Softfloat_roundingMode = mode
	t.Cleanup(func() { Softfloat_roundingMode = Softfloat_round_near_even })
}

// exact computes at 1000 bits in mode m. The mode only matters for the sign of an exact zero
// sum, which IEEE makes -0 when rounding down; every other result here is exact or is rounded
// again by roundTo in the same mode.
func exact(op string, a, b, c float64, m big.RoundingMode) *big.Float {
	x := new(big.Float).SetMode(m).SetPrec(1000).SetFloat64(a)
	y := new(big.Float).SetPrec(1000).SetFloat64(b)
	switch op {
	case "add":
		return x.Add(x, y)
	case "sub":
		return x.Sub(x, y)
	case "mul":
		return x.Mul(x, y)
	case "div":
		return x.Quo(x, y)
	case "sqrt":
		// big.Float.Sqrt is not correctly rounded in the directed modes: it can return the
		// root of a perfect square one 1000-bit ULP low. Rounded to nearest it is exact for
		// perfect squares, and otherwise far from any 24- or 53-bit boundary.
		return x.SetMode(big.ToNearestEven).Sqrt(x)
	}
	x.Mul(x, y)
	return x.Add(x, new(big.Float).SetPrec(1000).SetFloat64(c)) // mulAdd
}

// roundTo rounds v to prec bits in mode m. It reports false when the result would be outside
// the normal range, where math/big and IEEE disagree.
func roundTo(v *big.Float, prec uint, m big.RoundingMode, minExp, maxExp int) (*big.Float, bool) {
	r := new(big.Float).SetMode(m).SetPrec(prec).Set(v)
	if r.Sign() == 0 {
		return r, true
	}
	e := r.MantExp(nil) - 1 // r = 1.f × 2^e
	return r, e >= minExp && e <= maxExp
}

// normalOperand stays well inside the exponent range, so products and quotients mostly do too.
func normalOperand(r *rand.Rand, spread int) float64 {
	for {
		v := r.NormFloat64() * math.Pow(2, float64(r.Intn(2*spread)-spread))
		if r.Intn(4) == 0 {
			v = math.Round(v*16) / 16 // short mantissas make ties and exact results likely
		}
		if v != 0 {
			return v
		}
	}
}

func TestRoundingModes(t *testing.T) {
	const samples = 20_000
	ops := []string{"add", "sub", "mul", "div", "sqrt", "mulAdd"}
	for _, mode := range roundingModes {
		for _, op := range ops {
			t.Run("F64_"+op+"/"+mode.name, func(t *testing.T) {
				setRounding(t, mode.soft)
				r := rand.New(rand.NewSource(2))
				for i := 0; i < samples; i++ {
					a, b, c := normalOperand(r, 200), normalOperand(r, 200), normalOperand(r, 200)
					if op == "sqrt" {
						a = math.Abs(a)
					}
					want, ok := roundTo(exact(op, a, b, c, mode.big), 53, mode.big, -1022, 1023)
					if !ok {
						continue
					}
					wf, _ := want.Float64()
					var got Float64_t
					A, B, C := Float64_t(math.Float64bits(a)), Float64_t(math.Float64bits(b)), Float64_t(math.Float64bits(c))
					switch op {
					case "add":
						got = comp.F64_add(A, B)
					case "sub":
						got = comp.F64_sub(A, B)
					case "mul":
						got = comp.F64_mul(A, B)
					case "div":
						got = comp.F64_div(A, B)
					case "sqrt":
						got = comp.F64_sqrt(A)
					default:
						got = comp.F64_mulAdd(A, B, C)
					}
					if uint64(got) != math.Float64bits(wf) {
						t.Fatalf("F64_%s(%v, %v, %v) in %s = %#x (%v), want %#x (%v)", op, a, b, c,
							mode.name, uint64(got), math.Float64frombits(uint64(got)), math.Float64bits(wf), wf)
					}
				}
			})
			t.Run("F32_"+op+"/"+mode.name, func(t *testing.T) {
				setRounding(t, mode.soft)
				r := rand.New(rand.NewSource(3))
				for i := 0; i < samples; i++ {
					a := float64(float32(normalOperand(r, 30)))
					b := float64(float32(normalOperand(r, 30)))
					c := float64(float32(normalOperand(r, 30)))
					if op == "sqrt" {
						a = math.Abs(a)
					}
					want, ok := roundTo(exact(op, a, b, c, mode.big), 24, mode.big, -126, 127)
					if !ok {
						continue
					}
					wf, _ := want.Float32()
					var got Float32_t
					A := Float32_t(math.Float32bits(float32(a)))
					B := Float32_t(math.Float32bits(float32(b)))
					C := Float32_t(math.Float32bits(float32(c)))
					switch op {
					case "add":
						got = comp.F32_add(A, B)
					case "sub":
						got = comp.F32_sub(A, B)
					case "mul":
						got = comp.F32_mul(A, B)
					case "div":
						got = comp.F32_div(A, B)
					case "sqrt":
						got = comp.F32_sqrt(A)
					default:
						got = comp.F32_mulAdd(A, B, C)
					}
					if uint32(got) != math.Float32bits(wf) {
						t.Fatalf("F32_%s(%v, %v, %v) in %s = %#x (%v), want %#x (%v)", op, a, b, c,
							mode.name, uint32(got), math.Float32frombits(uint32(got)), math.Float32bits(wf), wf)
					}
				}
			})
		}
	}
}

// An exact zero sum of operands with opposite signs is +0 in every mode but round-down,
// where it is -0 (IEEE 754-2019, 6.3). That covers subtraction, addition of opposite signs,
// and a fused multiply-add whose product cancels the addend.
func TestExactZeroSign(t *testing.T) {
	one, negOne, two := f64(1), f64(-1), f64(2)
	one32, negOne32, two32 := f32(1), f32(-1), f32(2)
	for _, mode := range roundingModes {
		setRounding(t, mode.soft)
		want64, want32 := uint64(0), uint32(0)
		if mode.soft == Softfloat_round_min {
			want64, want32 = 1<<63, 1<<31
		}
		for name, got := range map[string]uint64{
			"1 - 1":      uint64(comp.F64_sub(one, one)),
			"1 + -1":     uint64(comp.F64_add(one, negOne)),
			"1 × 2 + -2": uint64(comp.F64_mulAdd(one, two, f64(-2))),
			"0 × 1 + -0": uint64(comp.F64_mulAdd(f64(0), one, f64(math.Copysign(0, -1)))),
		} {
			if got != want64 {
				t.Errorf("F64 %s in %s = %#x, want %#x", name, mode.name, got, want64)
			}
		}
		for name, got := range map[string]uint32{
			"1 - 1":      uint32(comp.F32_sub(one32, one32)),
			"1 + -1":     uint32(comp.F32_add(one32, negOne32)),
			"1 × 2 + -2": uint32(comp.F32_mulAdd(one32, two32, f32(-2))),
			"0 × 1 + -0": uint32(comp.F32_mulAdd(f32(0), one32, f32(float32(math.Copysign(0, -1))))),
		} {
			if got != want32 {
				t.Errorf("F32 %s in %s = %#x, want %#x", name, mode.name, got, want32)
			}
		}
		// Like signs never cancel: -0 + -0 is -0 in every mode.
		if got := uint64(comp.F64_add(f64(math.Copysign(0, -1)), f64(math.Copysign(0, -1)))); got != 1<<63 {
			t.Errorf("-0 + -0 in %s = %#x, want -0", mode.name, got)
		}
	}
}
