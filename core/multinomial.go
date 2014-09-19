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

func (m *Multinomial) Get(v string) *big.Rat {
	return m.Hist[v]
}

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

func (m *Multinomial) Accumulate(a *Multinomial) {
	for v, x := range a.Hist {
		acc(m.Hist[v], x)
	}
}
