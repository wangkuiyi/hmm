package core

import (
	"log"
	"math"
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
				m.S1[state] += 1
				m.S1Sum += 1
			} else { // Not the first one
				m.Σξ[prevState][state] += 1
			}
			if t < inst.T()-1 { // Not the last one.
				m.Σγ[state] += 1
			}
			for c := 0; c < C; c++ {
				for k, v := range inst.Obs[inst.index[t]][c] {
					m.Σγo[state][c].Inc(k, float64(v))
				}
			}
			prevState = state
		}
	}
	return m
}

func Train(corpus []*Instance, N, C, Iter int, baseline *Model) (*Model,
	[]float64) {
	var estimate *Model
	ll := make([]float64, 0)

	for iter := 0; iter < Iter; iter++ {
		estimate = NewModel(N, C)
		for _, inst := range corpus {
			β := β(inst, baseline)
			γ1, Σγ, Σξ, Σγo := Inference(inst, baseline, β)
			estimate.Update(γ1, Σγ, Σξ, Σγo)
		}

		for _, inst := range corpus {
			ll = append(ll, math.Log(Likelihood(inst, estimate)))
		}
		baseline = estimate
	}

	return estimate, ll
}

func β(inst *Instance, m *Model) [][]float64 {
	β := matrix(inst.T(), m.N())

	for t := inst.T() - 1; t >= 0; t-- {
		if t == inst.T()-1 {
			for i := 0; i < m.N(); i++ {
				β[t][i] = 1
			}
		} else {
			for i := 0; i < m.N(); i++ {
				sum := 0.0
				for j := 0; j < m.N(); j++ {
					sum += m.A(i, j) * m.B(j, inst.O(t+1)) * β[t+1][j]
				}
				β[t][i] = sum
			}
		}
	}

	return β
}

func αGen(inst *Instance, m *Model) func() []float64 {
	t := 0
	α := vector(m.N())
	return func() []float64 {
		if t == 0 { // Initialization
			for i := 0; i < m.N(); i++ {
				α[i] = m.π(i) * m.B(i, inst.O(0))
			}
		} else { // Induction
			nα := vector(m.N())
			for j := 0; j < m.N(); j++ {
				sum := 0.0
				for i := 0; i < m.N(); i++ {
					sum += α[i] * m.A(i, j)
				}
				nα[j] = sum * m.B(j, inst.O(t))
			}
			α = nα
		}
		t++
		return α
	}
}

func Inference(inst *Instance, m *Model, β [][]float64) (
	[]float64, []float64, [][]float64, [][]*Multinomial) {

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
		norm := 0.0
		for i := 0; i < m.N(); i++ {
			γ[i] = α[i] * β[t][i]
			norm += γ[i]
		}
		if norm != 0 {
			for i := 0; i < m.N(); i++ {
				γ[i] = γ[i] / norm
			}
		}

		// Accumulate γ(t) to γ1, Σγ, and Σγo.
		for i := 0; i < m.N(); i++ {
			if t == 0 {
				γ1[i] = γ[i]
			}

			if t < inst.T()-1 {
				Σγ[i] += γ[i]
			}

			for c := 0; c < m.C(); c++ {
				for k, v := range inst.O(t)[c] {
					Σγo[i][c].Inc(k, γ[i]*float64(v))
				}
			}
		}

		// Compute ξ and accumulate to Σξ.
		if t < inst.T()-1 {
			ξSum := 0.0
			for i := 0; i < m.N(); i++ {
				for j := 0; j < m.N(); j++ {
					x := α[i] * m.A(i, j) * m.B(j, inst.O(t+1)) * β[t+1][j]
					ξ[i][j] = x
					ξSum += x
				}
			}

			if ξSum != 0 {
				for i := 0; i < m.N(); i++ {
					for j := 0; j < m.N(); j++ {
						Σξ[i][j] += ξ[i][j] / ξSum
					}
				}
			}
		}
	}

	return γ1, Σγ, Σξ, Σγo
}

func Likelihood(inst *Instance, m *Model) float64 {
	gen := αGen(inst, m)
	for t := 0; t < inst.T(); t++ {
		α := gen()
		if t == inst.T()-1 {
			sum := 0.0
			for _, v := range α {
				sum += v
			}
			return sum
		}
	}
	return math.NaN()
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
