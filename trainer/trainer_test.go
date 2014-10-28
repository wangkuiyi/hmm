package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

const (
	kTrain = `github.com/wangkuiyi/hmm/trainer`
)

func TestTrain(t *testing.T) {
	goCompiler, e := exec.LookPath("go")
	if e != nil {
		t.Fatalf("Cannot find go in PATH: %v", e)
	}

	goPath := os.Getenv("GOPATH")
	if len(goPath) <= 0 {
		t.Fatalf("GOPATH not set")
	}

	if e := exec.Command(goCompiler, "install", kTrain).Run(); e != nil {
		t.Fatalf("Cannot build %s: %v", kTrain, e)
	}

	dir, e := ioutil.TempDir("", "wangkuiyi-hmm-trainer")
	if e != nil {
		t.Fatalf("Cannot create temp dir")
	}
	defer os.RemoveAll(dir)

	corpus := path.Join(goPath, "src", kTrain, "testdata/corpus.json")
	model := path.Join(dir, "model")
	logl := path.Join(dir, "logl")

	o, e := exec.Command(path.Join(goPath, "bin", path.Base(kTrain)),
		"-corpus="+corpus, "-model="+model, "-logl="+logl).CombinedOutput()
	if e != nil {
		t.Fatalf("%s failed: %s, %v", path.Base(kTrain), o, e)
	}

	if b, e := ioutil.ReadFile(model); e != nil {
		t.Fatalf("Failed reading model file: %s", model)
	} else if string(b) != truthModel {
		t.Errorf("Expecting %s, got %s", truthModel, string(b))
	}
}

var (
	truthModel = `{
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
)
