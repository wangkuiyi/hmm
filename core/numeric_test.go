package core

import (
	"math/big"
	"testing"
)

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
