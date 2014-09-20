package core

import (
	"log"
	"math/big"
)

// Denote the number of hiddden states by N, and the number of kinds
// of multinomial observables by C, the model is comprised of the
// following additive members:
//
//  S1[i]: the number of times that state i at the beginning of instances.
//  Σγ[i]: the expected number of state i.
//	Σξ[i][j]: the expected number of transitions from i to j.
//	Σγo[i][c][v]: the expected number of state i with v observed.
//
// We can derived the transition probabilities and observe
// probabilities as:
//
//  π[i] = S1[i]/sum(S1)
//  a[i][j] = Σξ[i][j] / Σγ[i]
//  b[i][c][v] =  Σγ[i][c][v] / Σγ[i][c].Sum
//
type Model struct {
	S1    []*big.Rat // Size is N
	S1Sum *big.Rat
	Σγ    []*big.Rat       // Size is N
	Σξ    [][]*big.Rat     // Size is N^2
	Σγo   [][]*Multinomial // Size is N*C
}

func NewModel(N, C int) *Model {
	if N <= 0 || C <= 0 {
		log.Panicf("Either is 0: N=%d, C=%d", N, C)
		return nil
	}

	return &Model{
		S1:    vector(N),
		S1Sum: zero(),
		Σγ:    vector(N),
		Σξ:    matrix(N, N),
		Σγo:   multinomialMatrix(N, C)}
}

func multinomialMatrix(x, y int) [][]*Multinomial {
	ret := make([][]*Multinomial, x)
	for i, _ := range ret {
		ret[i] = make([]*Multinomial, y)
		for j, _ := range ret[i] {
			ret[i][j] = NewMultinomial()
		}
	}
	return ret
}

func (m *Model) π(i int) *big.Rat {
	return div(m.S1[i], m.S1Sum)
}

func (m *Model) A(i, j int) *big.Rat {
	return div(m.Σξ[i][j], m.Σγ[i])
}

func (m *Model) B(state int, obs []Observed) *big.Rat {
	b := one()
	for c, ob := range obs {
		b.Mul(b, m.Σγo[state][c].Likelihood(ob))
	}
	return b
}

func (m *Model) Update(γ [][]*big.Rat, ξ [][]*Multinomial) {
	// TODO(wyi): implement it.
}

func (m *Model) N() int {
	return len(m.S1)
}

func (m *Model) C() int {
	return len(m.Σγo[0])
}
