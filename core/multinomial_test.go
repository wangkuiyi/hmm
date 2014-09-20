package core

import (
	"math/big"
	"testing"
)

func TestFact(t *testing.T) {
	if r := fact(0); !equ(r, one()) {
		t.Errorf("Expecting 1, got %d", r)
	}

	if r := fact(1); !equ(r, one()) {
		t.Errorf("Expecting 1, got %d", r)
	}

	if r := fact(2); !equ(r, big.NewRat(2, 1)) {
		t.Errorf("Expecting 2, got %d", r)
	}

	if r := fact(2); !equ(r, big.NewRat(2, 1)) {
		t.Errorf("Expecting 2, got %d", r)
	}

	truth := big.NewRat(2432902008176640000, 1)
	if r := fact(20); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}

func TestLikelihood(t *testing.T) {
	m := NewMultinomial()
	m.Acc("apple", one())
	m.Acc("orange", one())

	truth := big.NewRat(1, 2)
	if r := m.Likelihood(Observed{"apple": 1}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = big.NewRat(1, 4)
	if r := m.Likelihood(Observed{"orange": 2}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = zero()
	if r := m.Likelihood(Observed{"unknown": 2}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = big.NewRat(1, 2)
	if r := m.Likelihood(Observed{"apple": 1, "orange": 1}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = zero()
	if r := m.Likelihood(Observed{"apple": 1, "unknown": 1}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = big.NewRat(3, 8)
	if r := m.Likelihood(Observed{"apple": 2, "orange": 1}); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}
