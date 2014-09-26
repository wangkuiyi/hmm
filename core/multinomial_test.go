package core

import (
	"testing"
)

func TestFact(t *testing.T) {
	if r := fact(0); r != 1 {
		t.Errorf("Expecting 1, got %f", r)
	}

	if r := fact(1); r != 1 {
		t.Errorf("Expecting 1, got %f", r)
	}

	if r := fact(2); r != 2 {
		t.Errorf("Expecting 2, got %f", r)
	}

	if r := fact(2); r != 2 {
		t.Errorf("Expecting 2, got %f", r)
	}

	truth := 2432902008176640000.0
	if r := fact(20); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}

func TestAcc(t *testing.T) {
	m := NewMultinomial()
	if r := m.θ("apple"); r != 0 {
		t.Errorf("Expecting %v, got %v", 0, r)
	}

	m.Inc("apple", 10.0)
	if r := m.θ("apple"); r != 1 {
		t.Errorf("Expecting %v, got %v", 1, r)
	}

	m.Inc("orange", 5)
	truth := 2.0 / 3.0
	if r := m.θ("apple"); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
	truth = 1.0 / 3.0
	if r := m.θ("orange"); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}

func TestLikelihood(t *testing.T) {
	m := NewMultinomial()
	m.Inc("apple", 1)
	m.Inc("orange", 1)

	truth := 1.0 / 2.0
	if r := m.Likelihood(Observed{"apple": 1}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 1.0 / 4.0
	if r := m.Likelihood(Observed{"orange": 2}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 0
	if r := m.Likelihood(Observed{"unknown": 2}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 1.0 / 2.0
	if r := m.Likelihood(Observed{"apple": 1, "orange": 1}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 0
	if r := m.Likelihood(Observed{"apple": 1, "unknown": 1}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = 3.0 / 8.0
	if r := m.Likelihood(Observed{"apple": 2, "orange": 1}); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}
