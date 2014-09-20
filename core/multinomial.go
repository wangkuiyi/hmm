package core

import (
	"math/big"
)

// Multinomial represents a multinomial distribution -- the
// probability of observable "A" is Hist["A"]/Sum.
type Multinomial struct {
	Hist map[string]*big.Rat
	Sum  *big.Rat
}

func NewMultinomial() *Multinomial {
	return &Multinomial{
		Hist: make(map[string]*big.Rat),
		Sum:  zero()}
}

/*
func (m *Multinomial) Get(v string) *big.Rat {
	return m.Hist[v]
}
*/

func (m *Multinomial) Acc(v string, x *big.Rat) {
	if _, ok := m.Hist[v]; !ok {
		m.Hist[v] = zero()
	}
	acc(m.Hist[v], x)
	acc(m.Sum, x)
}

func (m *Multinomial) Inc(v string, x int) {
	if _, ok := m.Hist[v]; !ok {
		m.Hist[v] = zero()
	}
	inc(m.Hist[v], x)
	inc(m.Sum, x)
}

/*
func (m *Multinomial) Accumulate(a *Multinomial) {
	for v, x := range a.Hist {
		m.Acc(v, x)
	}
}
*/

func (m *Multinomial) Likelihood(ob Observed) *big.Rat {
	l := one()
	n := 0
	for k, c := range ob {
		l.Mul(l, pow(m.θ(k), c))
		l = div(l, fact(c))
		n += c
	}
	l.Mul(l, fact(n))
	return l
}

func (m *Multinomial) θ(key string) *big.Rat {
	if numerator, ok := m.Hist[key]; ok {
		return div(numerator, m.Sum)
	}
	return zero()
}

var (
	factorials []int64
)

func fact(x int) *big.Rat {
	if factorials == nil {
		factorials = make([]int64, 100)
		factorials[0] = 1
		factorials[1] = 1
		for i := int64(2); i < 100; i++ {
			factorials[i] = factorials[i-1] * i
		}
	}

	if x < 100 {
		return big.NewRat(factorials[x], 1)
	}
	f := factorials[99]
	for i := int64(100); i <= int64(x); i++ {
		f *= i
	}
	return big.NewRat(f, 1)
}
