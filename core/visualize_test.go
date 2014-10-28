package core

import (
	"io/ioutil"
	"testing"
)

func TestVisualizerOutputDot(t *testing.T) {
	e := NewVisualizer(kTruthModel).OutputDot("/tmp/a.dot", 1, 1)
	if e != nil {
		t.Fatalf("Visualizer.OutputDot failed: %v", e)
	}

	truth := `digraph Model {
node[shape=box,style="rounded,filled",fillcolor=azure];
start;
start -> 00001 [label="100.0%",weight=1,penwidth=3];
00000 [label="orange:100.0% \n"];
00001 [label="apple:100.0% \n"];
00000 -> 00001 [label="100.0%",weight=0,penwidth=3];
00001 -> 00000 [label="100.0%",weight=0,penwidth=3];
}
`
	if b, e := ioutil.ReadFile("/tmp/a.dot"); e != nil {
		t.Fatalf("Failed reading /tmp/a.dot")
	} else if string(b) != truth {
		t.Errorf("Expecting %s, got %s", truth, string(b))
	}
}
