package include

// for multi-thread version, you should keep the context for each thread

// Software floating-point underflow tininess-detection mode.
const (
	SoftfloatTininessBeforeRounding uint8 = 0
	SoftfloatTininessAfterRounding  uint8 = 1 // recommended default
)

// Software floating-point rounding mode.
const (
	Softfloat_round_near_even   uint8 = 0 // default
	Softfloat_round_minMag      uint8 = 1
	Softfloat_round_min         uint8 = 2
	Softfloat_round_max         uint8 = 3
	Softfloat_round_near_maxMag uint8 = 4
)

// Software floating-point exception flags.
const (
	Softfloat_flag_inexact   uint8 = 1
	Softfloat_flag_underflow uint8 = 2
	Softfloat_flag_overflow  uint8 = 4
	Softfloat_flag_infinite  uint8 = 8
	Softfloat_flag_invalid   uint8 = 16
)
