package core

import (
	"log"
	"math/big"
)

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

func β(inst *Instance, m *Model) [][]*big.Rat {
	/*
				β := createRatMatrix(inst.T(), m.N())

				for t := inst.T() - 1; t >= 0; t-- {
					if t == inst.T()-1 {
						for i := 0; i < m.N(); i++ {
							β[t][i] = one()
						}
					} else {
						for i := 0; i < m.N(); i++ {
							sum := zero()
							for j := 0; j < m.N(); j++ {
								acc(sum, prod(m.a(i, j),
		m.b(j, inst.Obs[inst.index[t]]), β[t+1][j]))
							}
							β[t][i] = sum
						}
					}
				}

				return β
	*/
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
