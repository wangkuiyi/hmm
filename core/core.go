package core

import (
	"log"
	"math/rand"
)

// Multinomial represents a multinomial distribution -- the
// probability of observable "A" is Hist["A"]/Sum.
type Multinomial struct {
	Hist map[string]float64
	Sum  float64
}

func NewMultinomial() *Multinomial {
	return new(Multinomial)
}

func (m *Multinomial) Get(v string) float64 {
	return m.Hist[v]
}

func (m *Multinomial) Add(v string, x float64) {
	m.Hist[v] += x
	m.Sum += x
}

func (m *Multinomial) Accumulate(a *Multinomial) {
	for v, x := range a.Hist {
		m.Hist[v] += x
	}
}

// Denote the number of hiddden states by K, and the number of kinds
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
	s1  []float64        // Size is K
	Σγ  []float64        // Size is K
	Σξ  [][]float64      // Size is K^2
	Σγo [][]*Multinomial // Size is K*C
}

func NewModel(K, C int) *Model {
	if K == 0 || C == 0 {
		log.Panicf("Either is 0: K=%d, C=%d", K, C)
		return nil
	}

	return &Model{
		s1:  make([]float64, K),
		Σγ:  make([]float64, K),
		Σξ:  makeMatrix(K, K),
		Σγo: makeMultinomialMatrix(K, C)}
}

func makeMatrix(x, y int) [][]float64 {
	ret := make([][]float64, x)
	for i, _ := range ret {
		ret[i] = make([]float64, y)
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

func (m *Model) Init(rng *rand.Rand) *Model {
	// TODO(wyi): implement it.
	return m
}

func (m *Model) Update(γ [][]float64, ξ [][]*Multinomial) {
	// TODO(wyi): implement it.
}

type Instance struct {
	Obs     [][]Observed // A sequence of observed.
	Periods []int        // Periods of above observed.

	index []int // Map time t to an observed.
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

func Train(corpus []*Instance, K, C, Iter int) *Model {
	rng := rand.New(rand.NewSource(0))
	baseline := NewModel(K, C).Init(rng)
	var estimate *Model

	for iter := 0; iter < Iter; iter++ {
		estimate = NewModel(K, C)
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

func β(inst *Instance, model *Model) [][]float64 {
	// TODO(wyi): implement it.
	return nil
}

func γ(inst *Instance, model *Model, β [][]float64) [][]float64 {
	// TODO(wyi): implement it.
	return nil
}

func ξ(inst *Instance, model *Model, β [][]float64) [][]*Multinomial {
	// TODO(wyi): implement it.
	return nil
}

func EstimateC(corpus []*Instance) int {
	// TODO(wyi): implement it.
	return 0
}
