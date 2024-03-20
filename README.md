# softfloat-go

This repository is an implementation of softfloat in Go, based on [berkeley-softfloat-3](https://github.com/ucb-bar/berkeley-softfloat-3).

## NOTE

1. If the result is NaN, return 0x7fc00000
2. int -0 = 10000000000000000000000000000000