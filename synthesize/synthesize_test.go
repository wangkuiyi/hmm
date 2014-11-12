package main

import (
	"github.com/wangkuiyi/buildnrun"
	"path"
	"testing"
)

func TestSynthesize(t *testing.T) {
	truth := `{"Obs":[[{"senior_staff":1}]],"Periods":[5]}
`
	pkg := "github.com/wangkuiyi/hmm/synthesize"
	model := "-model=" + path.Join(buildnrun.PkgDir(pkg), "testdata/ground_truth_model.json")
	if out, err, e := buildnrun.Run(pkg, "-instances=1", "-cardi=1", "-length=5", "-seed=0", model); e != nil {
		t.Errorf("Failed build and run %s: %v\nOutput:\n%s", pkg, e, err)
	} else if string(out) != truth {
		t.Errorf("Expecting %s\ngot %s", truth, out)
	}
}
