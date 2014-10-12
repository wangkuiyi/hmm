package core

import (
	"testing"
)

func TestVisualizerOutputDot(t *testing.T) {
	NewVisualizer(kTruthModel).OutputDot("/tmp/a.dot")
}
