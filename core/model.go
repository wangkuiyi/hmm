package core

import (
	"log"
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
	S1    []float64        // Size is N
	S1Sum float64          // sum_i S1[i]
	Σγ    []float64        // Size is N
	Σξ    [][]float64      // Size is N^2
	Σγo   [][]*Multinomial // Size is N*C
}

func NewModel(N, C int) *Model {
	if N <= 1 {
		log.Panicf("N=%d, must be >= 2", N)
	}
	if C <= 0 {
		log.Panicf("C=%d, must be >= 1", C)
		return nil
	}
	return &Model{
		S1:    vector(N),
		S1Sum: 0.0,
		Σγ:    vector(N),
		Σξ:    matrix(N, N),
		Σγo:   multinomialMatrix(N, C)}
}

func (m *Model) π(i int) float64 {
	if m.S1Sum != 0.0 {
		return m.S1[i] / m.S1Sum
	}
	return 0.0
}

func (m *Model) A(i, j int) float64 {
	if m.Σγ[i] != 0.0 {
		return m.Σξ[i][j] / m.Σγ[i]
	}
	return 0.0
}

func (m *Model) B(state int, obs []Observed) float64 {
	b := 1.0
	for c, ob := range obs {
		b *= m.Σγo[state][c].Likelihood(ob)
	}
	return b
}

func (m *Model) Update(γ1 []float64, Σγ []float64, Σξ [][]float64,
	Σγo [][]*Multinomial) {

	if len(γ1) != m.N() {
		log.Panicf("len(γ1) (%d) != m.N() (%d)", len(γ1), m.N())
	}
	for i := 0; i < m.N(); i++ {
		m.S1[i] += γ1[i]
		m.S1Sum += γ1[i]
	}

	if len(Σγ) != m.N() {
		log.Panicf("len(Σγ) (%d) != m.N() (%d)", len(Σγ), m.N())
	}
	for i := 0; i < m.N(); i++ {
		m.Σγ[i] += Σγ[i]
	}

	if len(Σξ) != m.N() {
		log.Panicf("len(Σξ) (%d) != m.N() (%d)", len(Σξ), m.N())
	}
	for i := 0; i < m.N(); i++ {
		if len(Σξ[i]) != m.N() {
			log.Panicf("len(Σξ[i]) (%d) != m.N() (%d)", len(Σξ[i]), m.N())
		}
		for j := 0; j < m.N(); j++ {
			m.Σξ[i][j] += Σξ[i][j]
		}
	}

	if len(Σγo) != m.N() {
		log.Panicf("len(Σγo) (%d) != m.N() (%d)", len(Σγo), m.N())
	}
	for i := 0; i < m.N(); i++ {
		if len(Σγo[i]) != m.C() {
			log.Panicf(" len(Σγo[i]) (%d) != m.C() (%d)", len(Σγo[i]), m.C())
		}
		for c := 0; c < m.C(); c++ {
			m.Σγo[i][c].Acc(Σγo[i][c])
		}
	}
}

func (m *Model) N() int {
	return len(m.S1)
}

func (m *Model) C() int {
	return len(m.Σγo[0])
}
