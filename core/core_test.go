package core

import (
	"math/big"
	"reflect"
	"testing"
)

var (
	kDachengObs = [][]Observed{
		[]Observed{
			Observed{"vice": 1, "president": 1},
			Observed{"applied": 1, "predictive": 1}},
		[]Observed{
			Observed{"head": 1, "solution": 1, "strategy": 1},
			Observed{"opera": 1, "solutions": 1}},
		[]Observed{
			Observed{"founder": 1},
			Observed{"helping": 1}},
		[]Observed{
			Observed{"senior": 1, "manager": 1},
			Observed{"linkedin": 1}}}
	kDachengPeriods = []int{8, 0, 1, 1}

	kGuanObs = [][]Observed{
		[]Observed{
			Observed{"cofounder": 1},
			Observed{"scissorsfly": 1}},
		[]Observed{
			Observed{"sr": 1, "associate": 1, "business": 1, "analyst": 1},
			Observed{"linkedin": 1}}}
	kGuanPeriods = []int{1, 1}

	kYiObs = [][]Observed{
		[]Observed{
			Observed{"software": 1, "engineer": 1},
			Observed{"google": 1}},
		[]Observed{
			Observed{"engineering": 1, "director": 1},
			Observed{"tencent": 1, "ads": 1}},
		[]Observed{
			Observed{"data": 1, "scientist": 1},
			Observed{"linkedin": 1}}}
	kYiPeriods = []int{3, 4, 0}

	kN = 2
	kC = 2
)

func TestBuildInstanceIndex(t *testing.T) {
	index := buildInstanceIndex(kDachengPeriods)
	exp := 0
	for _, l := range kDachengPeriods {
		exp += l
	}
	if len(index) != exp {
		t.Errorf("Expecting %d, got %d", exp, len(index))
	}
	truth := []int{0, 0, 0, 0, 0, 0, 0, 0, 2, 3}
	if !reflect.DeepEqual(index, truth) {
		t.Errorf("Expecting %d, got %d", truth, index)
	}
}

func TestNewInstance(t *testing.T) {
	dacheng := NewInstance(kDachengObs, kDachengPeriods)
	if !reflect.DeepEqual(dacheng.Obs, kDachengObs) {
		t.Errorf("Expecting %v, got %v", kDachengObs, dacheng.Obs)
	}
	if !reflect.DeepEqual(dacheng.Periods, kDachengPeriods) {
		t.Errorf("Expecting %v, got %v", kDachengPeriods, dacheng.Periods)
	}

	null := NewInstance([][]Observed{nil, nil}, []int{0, 0})
	if null != nil {
		t.Errorf("Expecting nil, got %v", null)
	}
}

func TestEstimateC(t *testing.T) {
	corpus := []*Instance{
		NewInstance(kDachengObs, kDachengPeriods),
		NewInstance(kGuanObs, kGuanPeriods),
		NewInstance(kYiObs, kYiPeriods)}
	c := EstimateC(corpus)
	if c != kC {
		t.Errorf("Expecting %d, got %d", kC, c)
	}
}

type mockRng struct {
	History []int
}

func (rng *mockRng) Intn(n int) int {
	if len(rng.History) == 0 {
		rng.History = make([]int, 1, 100)
		rng.History[0] = 0
		return 0
	}
	p := rng.History[len(rng.History)-1]
	if p+1 >= n {
		p = 0
	} else {
		p = p + 1
	}
	rng.History = append(rng.History, p)
	return p
}

func TestInit(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	c := EstimateC(corpus)
	rng := new(mockRng)
	m := Init(kN, c, corpus, rng)

	truth := &Model{
		S1: []*big.Rat{rat(1), rat(0)},
		Σγ: []*big.Rat{rat(5), rat(4)},
		Σξ: [][]*big.Rat{
			{rat(0), rat(5)},
			{rat(4), rat(0)}},
		Σγo: [][]*Multinomial{
			[]*Multinomial{
				&Multinomial{
					Hist: map[string]*big.Rat{
						"founder":   rat(1),
						"president": rat(4),
						"vice":      rat(4)},
					Sum: rat(9)},
				&Multinomial{
					Hist: map[string]*big.Rat{
						"applied":    rat(4),
						"helping":    rat(1),
						"predictive": rat(4)},
					Sum: rat(9)}},
			[]*Multinomial{
				&Multinomial{
					Hist: map[string]*big.Rat{
						"manager":   rat(1),
						"president": rat(4),
						"senior":    rat(1),
						"vice":      rat(4)},
					Sum: rat(10)},
				&Multinomial{
					Hist: map[string]*big.Rat{
						"applied":    rat(4),
						"linkedin":   rat(1),
						"predictive": rat(4)},
					Sum: rat(9)}}}}

	if !reflect.DeepEqual(m, truth) {
		t.Errorf("Expecting %v, got %v", truth, m)
	}
}

// func TestTrain(t *testing.T) {
// 	corpus := []*Instance{}

// 	C := EstimateC(corpus)
// 	N := 2
// 	Iter := 100
// 	model := Train(corpus, N, C, Iter)
// 	fmt.Println(model)
// }

func TestProd(t *testing.T) {
	if r := prod(); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}

	if r := prod(big.NewRat(1, 3), big.NewRat(3, 1)); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}

	a := big.NewRat(1, 3)
	b := big.NewRat(3, 1)
	if r := prod(a, b); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}
	if !equ(a, big.NewRat(1, 3)) {
		t.Errorf("Expecting %v, got %v", big.NewRat(1, 3), a)
	}
	if !equ(b, big.NewRat(3, 1)) {
		t.Errorf("Expecting %v, got %v", big.NewRat(3, 1), a)
	}
}

func TestDiv(t *testing.T) {
	if r := div(big.NewRat(1, 3), big.NewRat(1, 3)); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}

	a := big.NewRat(1, 3)
	if r := div(a, a); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}
	if !equ(a, big.NewRat(1, 3)) {
		t.Errorf("Expecting %v, got %v", big.NewRat(1, 3), a)
	}
}
