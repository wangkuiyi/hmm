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

func TestAcc(t *testing.T) {
	m := NewMultinomial()
	if r := m.θ("apple"); !equ(r, zero()) {
		t.Errorf("Expecting %v, got %v", zero(), r)
	}

	m.Inc("apple", rat(10))
	if r := m.θ("apple"); !equ(r, one()) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}

	m.Inc("orange", rat(5))
	truth := big.NewRat(2, 3)
	if r := m.θ("apple"); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", one(), r)
	}
	truth = big.NewRat(1, 3)
	if r := m.θ("orange"); !equ(r, truth) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}

func TestLikelihood(t *testing.T) {
	m := NewMultinomial()
	m.Inc("apple", one())
	m.Inc("orange", one())

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
