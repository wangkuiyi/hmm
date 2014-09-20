package core

import (
	"math/big"
	"testing"
)

func TestModelA(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	rng := new(mockRng)
	m := Init(kN, EstimateC(corpus), corpus, rng)

	if r := m.A(0, 0); !equ(r, zero()) {
		t.Errorf("Expecting %v, got %v", zero(), r)
	}

	if r := m.A(0, 1); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}

	if r := m.A(1, 0); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}

	if r := m.A(1, 1); !equ(r, zero()) {
		t.Errorf("Expecting %v, got %v", zero(), r)
	}
}

func TestModelB(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	rng := new(mockRng)
	m := Init(kN, EstimateC(corpus), corpus, rng)

	truth := big.NewRat(1, 9)
	if r := m.B(0, []Observed{{"founder": 1}, {}}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = big.NewRat(1, 81)
	obs := []Observed{{"founder": 1}, {"helping": 1}}
	if r := m.B(0, obs); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = zero()
	obs = []Observed{{"founder": 1}, {"unknown": 1}}
	if r := m.B(0, obs); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}
