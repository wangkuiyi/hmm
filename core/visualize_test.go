package core

import (
	"io/ioutil"
	"testing"
)

func TestVisualizerOutputDot(t *testing.T) {
	if e := NewVisualizer(kTruthModel).OutputDot("/tmp/a.dot"); e != nil {
		t.Fatalf("Visualizer.OutputDot failed: %v", e)
	}

	truth := `digraph Model {
start [shape=box];
start -> 00001 [label="1.000000",weight=2];
00000 [shape=ellipse,label="orange:1.000000 \n"];
00001 [shape=ellipse,label="apple:1.000000 \n"];
00000 -> 00001 [label="1.000000",weight=0];
00001 -> 00000 [label="1.000000",weight=2];
}
`
	if b, e := ioutil.ReadFile("/tmp/a.dot"); e != nil {
		t.Fatalf("Failed reading /tmp/a.dot")
	} else if string(b) != truth {
		t.Errorf("Expecting %s, got %s", truth, string(b))
	}
}
