package core

import (
	"math/big"
)

func zero() *big.Rat {
	return big.NewRat(0, 1)
}

func one() *big.Rat {
	return big.NewRat(1, 1)
}

func rat(n int) *big.Rat {
	return big.NewRat(int64(n), 1)
}

func equ(a *big.Rat, b *big.Rat) bool {
	return a.Cmp(b) == 0
}

func acc(r *big.Rat, x *big.Rat) {
	r.Add(r, x)
}

func inc(r *big.Rat, x int) {
	acc(r, big.NewRat(int64(x), 1))
}

func prod(v ...*big.Rat) *big.Rat {
	ret := one()
	for _, x := range v {
		ret.Mul(ret, x)
	}
	return ret
}

func div(a, b *big.Rat) *big.Rat {
	return prod(a, zero().Inv(b))
}

func pow(a *big.Rat, p int) *big.Rat {
	if p == 0 {
		return one()
	}

	inv := false
	if p < 0 {
		p = -p
		inv = true
	}

	ret := one()
	for i := 0; i < p; i++ {
		ret.Mul(ret, a)
	}

	if inv {
		return ret.Inv(ret)
	}
	return ret
}

func createRatVector(x int) []*big.Rat {
	ret := make([]*big.Rat, x)
	for i, _ := range ret {
		ret[i] = zero()
	}
	return ret
}

func createRatMatrix(x, y int) [][]*big.Rat {
	ret := make([][]*big.Rat, x)
	for i, _ := range ret {
		ret[i] = make([]*big.Rat, y)
		for j, _ := range ret[i] {
			ret[i][j] = zero()
		}
	}
	return ret
}
