package core

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
)

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

func TestInit(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	rng := new(mockRng)
	m := Init(kN, EstimateC(corpus), corpus, rng)

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

func TestBackward(t *testing.T) {
	inst := NewInstance(kDachengObs, kDachengPeriods)
	corpus := []*Instance{inst}
	rng := new(mockRng)
	m := Init(kN, EstimateC(corpus), corpus, rng)
	β := β(inst, m)
	truth := `[
  [
    "2305843009213693952/86552130074731456014931640625",
    "0"
  ],
  [
    "0",
    "9007199254740992/42741792629497015316015625"
  ],
  [
    "8796093022208/6514524101432253515625",
    "0"
  ],
  [
    "0",
    "34359738368/3217048938978890625"
  ],
  [
    "33554432/490329056390625",
    "0"
  ],
  [
    "0",
    "131072/242137805625"
  ],
  [
    "128/36905625",
    "0"
  ],
  [
    "0",
    "1/36450"
  ],
  [
    "1/450",
    "0"
  ],
  [
    "1",
    "1"
  ]
]`
	if b, e := json.MarshalIndent(β, "", "  "); e == nil {
		if string(b) != truth {
			t.Errorf("Expecting\n%v\ngot\n%v\n", truth, string(b))
		}
	} else {
		t.Errorf("json.MarshalIndent failed")
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
