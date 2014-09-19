package core

import (
	"log"
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

func (m *Multinomial) Add(v string, x *big.Rat) {
	if _, ok := m.Hist[v]; !ok {
		m.Hist[v] = zero()
	}
	add(m.Hist[v], x)
	add(m.Sum, x)
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
		m.Hist[v].Add(x, m.Hist[v])
	}
}

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
	S1  []*big.Rat       // Size is N
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
		S1:  createRatVector(N),
		Σγ:  createRatVector(N),
		Σξ:  createRatMatrix(N, N),
		Σγo: createRatHistMatrix(N, C)}
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

func createRatHistMatrix(x, y int) [][]*Multinomial {
	ret := make([][]*Multinomial, x)
	for i, _ := range ret {
		ret[i] = make([]*Multinomial, y)
		for j, _ := range ret[i] {
			ret[i][j] = NewMultinomial()
		}
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

type Rng interface {
	Intn(int) int
}

func Init(N, C int, corpus []*Instance, rng Rng) *Model {
	m := NewModel(N, C)

	for _, inst := range corpus {
		prevState := -1
		for t := 0; t < inst.T(); t++ {
			state := rng.Intn(N)
			if t == 0 { // Is the first element.
				inc(m.S1[state], 1)
			} else { // Not the first one
				inc(m.Σξ[prevState][state], 1)
			}
			if t < inst.T()-1 { // Not the last one.
				inc(m.Σγ[state], 1)
			}
			for c := 0; c < C; c++ {
				for k, v := range inst.Obs[inst.index[t]][c] {
					m.Σγo[state][c].Inc(k, v)
				}
			}
			prevState = state
		}
	}
	return m
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
	c := 0
	for _, inst := range corpus {
		for _, o := range inst.Obs {
			if c == 0 {
				c = len(o)
			} else if c != len(o) {
				log.Panicf("c = %d, len(o) = %d", c, len(o))
			}
		}
	}
	return c
}

func add(r *big.Rat, x *big.Rat) {
	r.Add(r, x)
}

func inc(r *big.Rat, x int) {
	add(r, big.NewRat(int64(x), 1))
}

func zero() *big.Rat {
	return big.NewRat(0, 1)
}

func one() *big.Rat {
	return big.NewRat(1, 1)
}

func rat(n int) *big.Rat {
	return big.NewRat(int64(n), 1)
}
