package core

import (
	"fmt"
	"github.com/wangkuiyi/parallel"
	"io"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
)

type Rng interface {
	Intn(int) int
}

func Init(N, C int, corpus []*Instance, rng Rng) *Model {
	m := NewModel(N, C)

	// TODO(wyi): Make this multithreading.
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

func Epoch(corpus []*Instance, N, C int, baseline *Model) *Model {
	estimate := NewModel(N, C)
	workers := runtime.NumCPU() - 1

	aggrγ1 := make(chan []float64)
	aggrΣγ := make(chan []float64)
	aggrΣξ := make(chan [][]float64)
	aggrΣγo := make(chan [][]*Multinomial)

	// TODO(wyi): Use WaitGroup here.
	go func() {
		for γ1 := range aggrγ1 {
			estimate.updateγ1(γ1)
		}
	}()
	go func() {
		for Σγ := range aggrΣγ {
			estimate.updateΣγ(Σγ)
		}
	}()
	go func() {
		for Σξ := range aggrΣξ {
			estimate.updateΣξ(Σξ)
		}
	}()
	go func() {
		for Σγo := range aggrΣγo {
			estimate.updateΣγo(Σγo)
		}
	}()

	parallel.For(0, workers, 1, func(worker int) {
		for i := worker; i < len(corpus); i += workers {
			β := β(corpus[i], baseline)
			γ1, Σγ, Σξ, Σγo := Inference(corpus[i], baseline, β)
			aggrγ1 <- γ1
			aggrΣγ <- Σγ
			aggrΣξ <- Σξ
			aggrΣγo <- Σγo
		}
	})
	close(aggrγ1)
	close(aggrΣγ)
	close(aggrΣξ)
	close(aggrΣγo)

	return estimate
}

// LogL returns the per-instance log-likelihood.  Note that the log-likelihood
// of some instances might evaluate to NaN.  LogL will ignore such cases.
func LogL(corpus []*Instance, model *Model) float64 {
	logl := 0.0
	workers := runtime.NumCPU() - 1
	aggr := make(chan float64)
	inst := 0

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for l := range aggr {
			if !math.IsNaN(l) && l > 0 {
				logl += math.Log(l)
				inst++
			}
		}
		wg.Done()
	}()

	parallel.For(0, workers, 1, func(worker int) {
		for i := worker; i < len(corpus); i++ {
			aggr <- Likelihood(corpus[i], model)
		}
	})
	close(aggr)

	wg.Wait()
	if inst == 0 {
		log.Fatalf("All instances have log-likelihood evaluated to NaN.")
	}
	return logl / float64(inst)
}

func Train(corpus []*Instance, N, C, I int, m *Model, ll io.Writer) *Model {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	for iter := 0; iter < I; iter++ {
		select {
		case <-sig:
			log.Printf("Terminate due to signal")
			return m
		default:
			m = Epoch(corpus, N, C, m)
			// The output format is compatible with Gnuplot data file format.
			fmt.Fprintf(ll, "%d\t%f\n", iter, LogL(corpus, m))
		}
	}
	return m
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
