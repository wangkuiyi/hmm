package main

import (
	"github.com/wangkuiyi/buildnrun"
	"path"
	"testing"
)

func TestTrain(t *testing.T) {
	trainer := "github.com/wangkuiyi/hmm/trainer"
	corpus := buildnrun.Pkg(path.Join(trainer, "testdata/corpus.json"))
	truth := `{
  "S1": [
    0,
    2
  ],
  "S1Sum": 2,
  "Σγ": [
    4,
    6
  ],
  "Σξ": [
    [
      0,
      4
    ],
    [
      6,
      0
    ]
  ],
  "Σγo": [
    [
      {
        "Hist": {
          "orange": 6
        },
        "Sum": 6
      }
    ],
    [
      {
        "Hist": {
          "apple": 6
        },
        "Sum": 6
      }
    ]
  ]
}`
	if out, err, e := buildnrun.Run(trainer, "-corpus="+corpus, "-logl=/dev/null"); e != nil {
		t.Fatalf("Failed build and run trainer: %v, %s", e, err)
	} else if out != truth {
		t.Errorf("Expecting %s, got %s", truth, out)
	}
}
