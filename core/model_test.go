package core

import (
	"testing"
)

func TestModelA(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	rng := new(mockRng)
	m := Init(kN, EstimateC(corpus), corpus, rng)

	if r := m.A(0, 0); r != 0 {
		t.Errorf("Expecting %v, got %v", 0.0, r)
	}

	if r := m.A(0, 1); r != 1 {
		t.Errorf("Expecting %v, got %v", 1.0, r)
	}

	if r := m.A(1, 0); r != 1 {
		t.Errorf("Expecting %v, got %v", 1.0, r)
	}

	if r := m.A(1, 1); r != 0 {
		t.Errorf("Expecting %v, got %v", 0, r)
	}
}

func TestModelB(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	rng := new(mockRng)
	m := Init(kN, EstimateC(corpus), corpus, rng)

	truth := 1.0 / 9.0
	if r := m.B(0, []Observed{{"founder": 1}, {}}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 1.0 / 81.0
	obs := []Observed{{"founder": 1}, {"helping": 1}}
	if r := m.B(0, obs); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 0.0
	obs = []Observed{{"founder": 1}, {"unknown": 1}}
	if r := m.B(0, obs); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}
