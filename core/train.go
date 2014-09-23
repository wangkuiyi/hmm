package core

import (
	"fmt"
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
				inc(m.S1Sum, 1)
			} else { // Not the first one
				inc(m.Σξ[prevState][state], 1)
			}
			if t < inst.T()-1 { // Not the last one.
				inc(m.Σγ[state], 1)
			}
			for c := 0; c < C; c++ {
				for k, v := range inst.Obs[inst.index[t]][c] {
					m.Σγo[state][c].Inc(k, rat(v))
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
		fmt.Println("Iter ", iter)
		estimate = NewModel(N, C)
		for _, inst := range corpus {
			β := β(inst, baseline)
			γ1, Σγ, Σξ, Σγo := Inference(inst, baseline, β)
			estimate.Update(γ1, Σγ, Σξ, Σγo)
		}
		baseline = estimate
	}

	return estimate
}

func β(inst *Instance, m *Model) [][]*big.Rat {
	β := matrix(inst.T(), m.N())

	for t := inst.T() - 1; t >= 0; t-- {
		if t == inst.T()-1 {
			for i := 0; i < m.N(); i++ {
				β[t][i] = one()
			}
		} else {
			for i := 0; i < m.N(); i++ {
				sum := zero()
				for j := 0; j < m.N(); j++ {
					acc(sum, prod(m.A(i, j), m.B(j, inst.O(t+1)), β[t+1][j]))
				}
				β[t][i] = sum
			}
		}
	}

	return β
}

func αGen(inst *Instance, m *Model) func() []*big.Rat {
	t := 0
	α := vector(m.N())
	return func() []*big.Rat {
		if t == 0 { // Initialization
			for i := 0; i < m.N(); i++ {
				α[i] = prod(m.π(i), m.B(i, inst.O(0)))
			}
		} else { // Induction
			nα := vector(m.N())
			for j := 0; j < m.N(); j++ {
				sum := zero()
				for i := 0; i < m.N(); i++ {
					acc(sum, prod(α[i], m.A(i, j)))
				}
				nα[j] = prod(sum, m.B(j, inst.O(t)))
			}
			α = nα
		}
		t++
		return α
	}
}

func Inference(inst *Instance, m *Model, β [][]*big.Rat) (
	[]*big.Rat, []*big.Rat, [][]*big.Rat, [][]*Multinomial) {

	γ1 := vector(m.N())
	Σγ := vector(m.N())
	Σγo := multinomialMatrix(m.N(), m.C())
	Σξ := matrix(m.N(), m.N())

	gen := αGen(inst, m)
	γ := vector(m.N())
	ξ := matrix(m.N(), m.N())

	for t := 0; t < inst.T(); t++ {
		α := gen()

		// Compute γ(t).
		norm := zero()
		for i := 0; i < m.N(); i++ {
			γ[i] = prod(α[i], β[t][i])
			acc(norm, γ[i])
		}
		if !equ(norm, zero()) {
			for i := 0; i < m.N(); i++ {
				γ[i] = div(γ[i], norm)
			}
		}

		// Accumulate γ(t) to γ1, Σγ, and Σγo.
		for i := 0; i < m.N(); i++ {
			if t == 0 {
				γ1[i] = γ[i]
			}

			if t < inst.T()-1 {
				acc(Σγ[i], γ[i])
			}

			for c := 0; c < m.C(); c++ {
				for k, v := range inst.O(t)[c] {
					Σγo[i][c].Inc(k, prod(γ[i], rat(v)))
				}
			}
		}

		// Compute ξ and accumulate to Σξ.
		if t < inst.T()-1 {
			ξSum := zero()
			for i := 0; i < m.N(); i++ {
				for j := 0; j < m.N(); j++ {
					x := prod(α[i], m.A(i, j), m.B(j, inst.O(t+1)), β[t+1][j])
					ξ[i][j] = x
					acc(ξSum, x)
				}
			}

			if !equ(ξSum, zero()) {
				for i := 0; i < m.N(); i++ {
					for j := 0; j < m.N(); j++ {
						acc(Σξ[i][j], div(ξ[i][j], ξSum))
					}
				}
			}
		}
	}

	return γ1, Σγ, Σξ, Σγo
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
