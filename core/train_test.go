package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
)

var (
	kSimpleObs = [][]Observed{
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
		[]Observed{Observed{"apple": 1}},
		[]Observed{Observed{"orange": 1}},
	}
	kSimplePeriods = []int{1, 1, 1, 1, 1, 1}

	kTruthModel = &Model{
		S1:    []float64{0, 2},
		S1Sum: 2,
		Σγ:    []float64{4, 6},
		Σξ:    [][]float64{{0, 4}, {6, 0}},
		Σγo: [][]*Multinomial{
			{&Multinomial{
				Hist: map[string]float64{"orange": 6},
				Sum:  6,
			}},
			{&Multinomial{
				Hist: map[string]float64{"apple": 6},
				Sum:  6,
			}}}}
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
    1,
    0
  ],
  "S1Sum": 1,
  "Σγ": [
    5,
    4
  ],
  "Σξ": [
    [
      0,
      5
    ],
    [
      4,
      0
    ]
  ],
  "Σγo": [
    [
      {
        "Hist": {
          "founder": 1,
          "president": 4,
          "vice": 4
        },
        "Sum": 9
      },
      {
        "Hist": {
          "applied": 4,
          "helping": 1,
          "predictive": 4
        },
        "Sum": 9
      }
    ],
    [
      {
        "Hist": {
          "manager": 1,
          "president": 4,
          "senior": 1,
          "vice": 4
        },
        "Sum": 10
      },
      {
        "Hist": {
          "applied": 4,
          "linkedin": 1,
          "predictive": 4
        },
        "Sum": 9
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
    2.6641089101132097e-11,
    0
  ],
  [
    0,
    2.1073517746012688e-10
  ],
  [
    1.3502280266756764e-09,
    0
  ],
  [
    0,
    1.0680514664133765e-08
  ],
  [
    6.843247725720863e-08,
    0
  ],
  [
    0,
    5.413115876790916e-07
  ],
  [
    3.468305983166524e-06,
    0
  ],
  [
    0,
    2.7434842249657068e-05
  ],
  [
    0.0022222222222222227,
    0
  ],
  [
    1,
    1
  ]
]`
	if b, e := json.MarshalIndent(β, "", "  "); e == nil {
		if string(b) != truth {
			fmt.Printf("Expecting\n%v\ngot\n%v\n", truth, string(b))
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

	truth := "[0.1560737692424935 0]"
	if r := fmt.Sprint(αGen()); r != truth {
		t.Errorf("Expecting %v, got %v", truth, r)
	}

	truth = "[0 0.019730807370902888]"
	if r := fmt.Sprint(αGen()); r != truth {
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

func TestTrain(t *testing.T) {
	corpus := []*Instance{
		NewInstance(kSimpleObs, kSimplePeriods),
		NewInstance(kSimpleObs, kSimplePeriods)}
	C := EstimateC(corpus)
	N := 2
	Iter := 20

	baseline := Init(N, C, corpus, rand.New(rand.NewSource(99)))
	model, _ := Train(corpus, N, C, Iter, baseline)

	if eq, b1, b2, e := jsonEncodingEqu(model, kTruthModel); e != nil {
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

func TestLikelihood(t *testing.T) {
	if l := Likelihood(NewInstance(kSimpleObs, kSimplePeriods),
		kTruthModel); l != 1.0 {
		t.Errorf("Expecting 1, got %f", l)
	}
}
