package core

import (
	"math"
	"math/rand"
	"sort"
)

// Multinomial represents a multinomial distribution -- the
// probability of observable "A" is Hist["A"]/Sum.
type Multinomial struct {
	Hist map[string]float64
	Sum  float64
}

func NewMultinomial() *Multinomial {
	return &Multinomial{
		Hist: make(map[string]float64),
		Sum:  0.0}
}

func (m *Multinomial) Inc(v string, x float64) {
	if x != 0 {
		m.Hist[v] += x
		m.Sum += x
	}
}

func (m *Multinomial) Acc(a *Multinomial) {
	for v, x := range a.Hist {
		m.Inc(v, x)
	}
}

func (m *Multinomial) Likelihood(ob Observed) float64 {
	l := 1.0
	n := 0
	for k, c := range ob {
		l *= math.Pow(m.θ(k), float64(c))
		l /= fact(c)
		n += c
	}
	l *= fact(n)
	return l
}

func (m *Multinomial) θ(key string) float64 {
	if numerator, ok := m.Hist[key]; ok {
		return numerator / m.Sum
	}
	return 0
}

func (m *Multinomial) Sample(n int, rng *rand.Rand) map[string]int {
	ret := make(map[string]int)

	// We have to copy keys out from map m.Hist to a slice and sort
	// them. Otherwise, we would have to access m.Hist directly using
	// range, which introduces randomness and leads to random output
	// even given a fix-seeded rng.
	keys := make([]string, 0, len(m.Hist))
	for k, _ := range m.Hist {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i < n; i++ {
		p := rng.Float64() * m.Sum
		sum := 0.0
		for _, k := range keys {
			sum += m.Hist[k]
			if p < sum {
				ret[k]++
				break
			}
		}
	}
	return ret
}

var (
	factorials []int64
)

func fact(x int) float64 {
	if factorials == nil {
		factorials = make([]int64, 100)
		factorials[0] = 1
		factorials[1] = 1
		for i := int64(2); i < 100; i++ {
			factorials[i] = factorials[i-1] * i
		}
	}

	if x < 100 {
		return float64(factorials[x])
	}
	f := factorials[99]
	for i := int64(100); i <= int64(x); i++ {
		f *= i
	}
	return float64(f)
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
