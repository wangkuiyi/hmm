package core

import (
	"log"
	"math/big"
	"math/rand"
)

// Multinomial represents a multinomial distribution -- the
// probability of observable "A" is Hist["A"]/Sum.
type Multinomial struct {
	Hist map[string]*big.Rat
	Sum  *big.Rat
}

func NewMultinomial() *Multinomial {
	return new(Multinomial)
}

func (m *Multinomial) Get(v string) *big.Rat {
	return m.Hist[v]
}

func (m *Multinomial) Add(v string, x *big.Rat) {
	m.Hist[v].Add(x, m.Hist[v])
	m.Sum.Add(x, m.Sum)
}

func (m *Multinomial) Accumulate(a *Multinomial) {
	for v, x := range a.Hist {
		m.Hist[v].Add(x, m.Hist[v])
	}
}

// Denote the number of hiddden states by N, and the number of kinds
// of multinomial observables by C, the model is comprised of the
// following additive members:
//
//  s1[i]: the number of times that state i at the beginning of instances.
//  Σγ[i]: the expected number of state i.
//	Σξ[i][j]: the expected number of transitions from i to j.
//	Σγo[i][c][v]: the expected number of state i with v observed.
//
// We can derived the transition probabilities and observe probabilities as:
//
//  π[i] = s1[i]/sum(s1)
//  a[i][j] = Σξ[i][j] / Σγ[i]
//  b[i][c][v] =  Σγ[i][c][v] / Σγ[i][c].Sum
//
type Model struct {
	s1  []*big.Rat       // Size is N
	Σγ  []*big.Rat       // Size is N
	Σξ  [][]*big.Rat     // Size is N^2
	Σγo [][]*Multinomial // Size is N*C
}

func NewModel(N, C int) *Model {
	if N == 0 || C == 0 {
		log.Panicf("Either is 0: N=%d, C=%d", N, C)
		return nil
	}

	return &Model{
		s1:  make([]*big.Rat, N),
		Σγ:  make([]*big.Rat, N),
		Σξ:  makeMatrix(N, N),
		Σγo: makeMultinomialMatrix(N, C)}
}

func makeMatrix(x, y int) [][]*big.Rat {
	ret := make([][]*big.Rat, x)
	for i, _ := range ret {
		ret[i] = make([]*big.Rat, y)
	}
	return ret
}

func makeMultinomialMatrix(x, y int) [][]*Multinomial {
	ret := make([][]*Multinomial, x)
	for i, _ := range ret {
		ret[i] = make([]*Multinomial, y)
	}
	return ret
}

func (m *Model) Update(γ [][]*big.Rat, ξ [][]*Multinomial) {
	// TODO(wyi): implement it.
}

type Instance struct {
	Obs     [][]Observed // A sequence of observed.
	Periods []int        // Periods of above observed.
	index   []int        // Map time t to an observed.
}

type Observed map[string]int

func NewInstance(obs [][]Observed, periods []int) *Instance {
	if len(periods) != len(obs) {
		log.Panicf("len(period)=%d, len(obs)=%d", len(periods), len(obs))
	}

	ret := &Instance{
		Obs:     obs,
		Periods: periods,
		index:   buildInstanceIndex(periods)}

	if len(ret.index) == 0 {
		log.Printf("periods are all 0: %v, obs: %v", periods, obs)
		return nil
	}
	return ret
}

func buildInstanceIndex(periods []int) []int {
	var T int
	for _, l := range periods {
		T += l
	}
	if T <= 0 {
		return nil
	}

	ret := make([]int, 0, T)
	for i, l := range periods {
		for j := 0; j < l; j++ {
			ret = append(ret, i)
		}
	}
	return ret
}

func (i *Instance) T() int {
	return len(i.index)
}

func Init(N, C int, corpus []*Instance, rng *rand.Rand) *Model {
	model := NewModel(N, C)

	return model
}

func Train(corpus []*Instance, N, C, Iter int, baseline *Model) *Model {
	var estimate *Model

	for iter := 0; iter < Iter; iter++ {
		estimate = NewModel(N, C)
		for _, inst := range corpus {
			β := β(inst, baseline)
			γ := γ(inst, baseline, β)
			ξ := ξ(inst, baseline, β)
			estimate.Update(γ, ξ)
		}
		baseline = estimate
	}

	return estimate
}

func β(inst *Instance, model *Model) [][]*big.Rat {
	// TODO(wyi): implement it.
	return nil
}

func γ(inst *Instance, model *Model, β [][]*big.Rat) [][]*big.Rat {
	// TODO(wyi): implement it.
	return nil
}

func ξ(inst *Instance, model *Model, β [][]*big.Rat) [][]*Multinomial {
	// TODO(wyi): implement it.
	return nil
}

func EstimateC(corpus []*Instance) int {
	// TODO(wyi): implement it.
	return 0
}
