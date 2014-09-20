package core

import (
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
