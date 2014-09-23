package core

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
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
	m := Init(kN, EstimateC(corpus), corpus, new(mockRng))

	truth := `{
  "S1": [
    "1",
    "0"
  ],
  "S1Sum": "1",
  "Σγ": [
    "5",
    "4"
  ],
  "Σξ": [
    [
      "0",
      "5"
    ],
    [
      "4",
      "0"
    ]
  ],
  "Σγo": [
    [
      {
        "Hist": {
          "founder": "1",
          "president": "4",
          "vice": "4"
        },
        "Sum": "9"
      },
      {
        "Hist": {
          "applied": "4",
          "helping": "1",
          "predictive": "4"
        },
        "Sum": "9"
      }
    ],
    [
      {
        "Hist": {
          "manager": "1",
          "president": "4",
          "senior": "1",
          "vice": "4"
        },
        "Sum": "10"
      },
      {
        "Hist": {
          "applied": "4",
          "linkedin": "1",
          "predictive": "4"
        },
        "Sum": "9"
      }
    ]
  ]
}`
	if b, e := json.MarshalIndent(m, "", "  "); e != nil {
		t.Errorf("json.MarshalIndent failed: %v", e)
	} else {
		if string(b) != truth {
			t.Errorf("Expecting\n%s\ngot\n%s\n", truth, b)
		}
	}
}

func TestBackward(t *testing.T) {
	inst := NewInstance(kDachengObs, kDachengPeriods)
	corpus := []*Instance{inst}
	m := Init(kN, EstimateC(corpus), corpus, new(mockRng))

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

func TestForwardGenerator(t *testing.T) {
	inst := NewInstance(kDachengObs, kDachengPeriods)
	corpus := []*Instance{inst}
	m := Init(kN, EstimateC(corpus), corpus, new(mockRng))

	αGen := αGen(inst, m)

	truth := []*big.Rat{big.NewRat(1024, 6561), zero()}
	if r := αGen(); !equ(r[0], truth[0]) || !equ(r[1], truth[1]) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = []*big.Rat{zero(), big.NewRat(262144, 13286025)}
	if r := αGen(); !equ(r[0], truth[0]) || !equ(r[1], truth[1]) {
		t.Errorf("Expecting %v, got %v", truth, r)
	}
}

func TestInference(t *testing.T) {
	inst := NewInstance(kDachengObs, kDachengPeriods)
	corpus := []*Instance{inst}
	m := Init(kN, EstimateC(corpus), corpus, new(mockRng))

	// The HMM decoding of an instance used to train (initialize) the
	// HMM model should have statitics exactly match the model.
	β := β(inst, m)
	γ1, Σγ, Σξ, Σγo := Inference(inst, m, β)

	if eq, b1, b2, e := jsonEncodingEqu(γ1, m.S1); e != nil {
		t.Errorf("json.MarshalIndent: %v", e)
	} else if !eq {
		t.Errorf("Expecting\n%s\ngot\n%s\n", b2, b1)
	}

	if eq, b1, b2, e := jsonEncodingEqu(Σγ, m.Σγ); e != nil {
		t.Errorf("json.MarshalIndent: %v", e)
	} else if !eq {
		t.Errorf("Expecting\n%s\ngot\n%s\n", b2, b1)
	}

	if eq, b1, b2, e := jsonEncodingEqu(Σγo, m.Σγo); e != nil {
		t.Errorf("json.MarshalIndent: %v", e)
	} else if !eq {
		t.Errorf("Expecting\n%s\ngot\n%s\n", b2, b1)
	}

	if eq, b1, b2, e := jsonEncodingEqu(Σξ, m.Σξ); e != nil {
		t.Errorf("json.MarshalIndent: %v", e)
	} else if !eq {
		t.Errorf("Expecting\n%s\ngot\n%s\n", b2, b1)
	}
}

func jsonEncodingEqu(v1, v2 interface{}) (bool, []byte, []byte, error) {
	b1, e := json.MarshalIndent(v1, "", "  ")
	if e != nil {
		return false, nil, nil, fmt.Errorf("json.MarshalIndent: %v", e)
	}

	b2, e := json.MarshalIndent(v2, "", "  ")
	if e != nil {
		return false, nil, nil, fmt.Errorf("json.MarshalIndent: %v", e)
	}

	eq := string(b1) == string(b2)
	return eq, b1, b2, nil
}

func TestTrain(t *testing.T) {
	kSimpleObs := [][]Observed{
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}}}
	kSimplePeriods := []int{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	corpus := []*Instance{
		NewInstance(kSimpleObs, kSimplePeriods)}
	C := EstimateC(corpus)
	N := 2
	Iter := 20

	baseline := Init(N, C, corpus, rand.New(rand.NewSource(99)))
	model := Train(corpus, N, C, Iter, baseline)

	if b, e := json.MarshalIndent(model, "", "  "); e == nil {
		fmt.Printf("%s\n", b) // debug
	} else {
		t.Errorf("json.MarshalIndent: %v", e)
	}
}
