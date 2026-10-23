# softfloat-go

A Go port of [Berkeley SoftFloat 3e](https://github.com/ucb-bar/berkeley-softfloat-3):
IEEE 754 binary16, binary32 and binary64 arithmetic done entirely in software, so every host
computes the same bits. It follows RISC-V's choices for NaNs, FMIN/FMAX, FCLASS and integer
conversion.

It is the floating-point unit of [MLVM](https://github.com/yao-jz/mlvm), a deterministic
RISC-V virtual machine, and is kept in sync with MLVM's copy (`rvgo/comp`).

## Use

```go
import (
	"math"

	"github.com/yao-jz/softfloat-go/comp"
	sf "github.com/yao-jz/softfloat-go/include"
)

sf.Softfloat_roundingMode = sf.Softfloat_round_near_even
sf.Softfloat_exceptionFlags = 0
z := comp.F64_div(sf.Float64_t(math.Float64bits(1)), sf.Float64_t(math.Float64bits(3)))
// math.Float64frombits(uint64(z)) == 0.3333333333333333
// sf.Softfloat_exceptionFlags == sf.Softfloat_flag_inexact
```

Values are raw bit patterns: `Float16_t`, `Float32_t` and `Float64_t` are `uint16`,
`uint32` and `uint64`. The rounding mode and the exception flags are package-level variables,
as in the C library, so use the package from one goroutine at a time.

## Packages

- `comp`: add, sub, mul, div, rem, sqrt, fused mulAdd and roundToInt; quiet and signaling
  comparisons; conversions between f16, f32, f64 and 32- and 64-bit integers; RISC-V's
  FMIN/FMAX (`F32_min`, `F64_max`, ...) and FCLASS (`F32_classS`, `F64_classS`).
- `include`: the types, the five rounding modes and the five exception flags.
- `riscv`: the RISC-V specialization. NaN results are the canonical NaN (`0x7FC00000`,
  `0x7FF8000000000000`), a signaling NaN operand raises invalid, and out-of-range integer
  conversions saturate and raise invalid.

## Tests

```bash
go test ./...
```

- Every add, sub, mul, div, sqrt and fused multiply-add result is compared bit for bit with
  the host FPU: 200,000 operand sets per operation and width, skewed toward subnormals, ties,
  wide exponent gaps and special values.
- The four directed and away-from-zero rounding modes are checked against exact `math/big`
  results.
- Exception flags, quiet and signaling comparisons, FCLASS, FMIN/FMAX, and integer
  conversions, including RISC-V's saturation and the sign of exact zero results.

## Changes in October 2026

Synced from MLVM's copy, which fixed and extended this code:

- Rounding modes other than round-to-nearest-even, and exception flags. The 2024 version
  always rounded to nearest and raised no flags.
- `F64_div` was one ULP off on about 2% of inputs (a `uint32` truncation dropped from the C).
- `F64_to_ui64` ignored its operand and returned 0.
- A signaling NaN operand now raises invalid, and FEQ is quiet while FLT and FLE signal.
- New: FMIN/FMAX and FCLASS with RISC-V semantics.

Fixed here, found by the new tests:

- An exact zero from subtraction or fused multiply-add is -0 when rounding down, as IEEE 754
  requires; it was always +0.
- The package builds: files imported both `RISCV` and `riscv`, which Go rejects.

## License

MIT for the port; the Berkeley SoftFloat license applies to the code it derives from. Both are
in [LICENSE](LICENSE).
