package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
)

type Visualizer struct {
	*Model
}

func NewVisualizer(m *Model) *Visualizer {
	return &Visualizer{m}
}

func (v *Visualizer) Draw(filename string) error {
	dir, e := ioutil.TempDir("", "Visualizer")
	if e != nil {
		return fmt.Errorf("Cannot create temp dir: %v", e)
	}
	defer os.RemoveAll(dir)

	dot := path.Join(dir, "model.dot")
	if e := v.OutputDot(dot); e != nil {
		return fmt.Errorf("Failed output .dot file: %v", e)
	}

	b, e := exec.Command("dot", "-Tpdf", "-o", filename, dot).CombinedOutput()
	if e != nil {
		return fmt.Errorf("Failed execution dot: %s, %v", b, e)
	}

	return nil

}

func (v *Visualizer) OutputDot(dotfile string) error {
	f, e := os.Create(dotfile)
	if e != nil {
		return fmt.Errorf("Cannot create dot file: %v", e)
	}
	defer f.Close()

	fmt.Fprintf(f, "digraph Model { \n")
	v.formatInits(f)
	v.formatNodes(f)
	v.formatEdges(f, v.thresholdEdgeWeight(1))
	fmt.Fprintf(f, "}\n")
	return nil
}

func (v *Visualizer) formatInits(w io.Writer) {
	fmt.Fprintf(w, "start [shape=box]; \n")
	for i, p := range v.S1 {
		if p > 0 {
			fmt.Fprintf(w, "start -> %05d [label=\"%f\",weight=%d];\n",
				i, float64(p)/float64(v.S1Sum), int(p))
		}
	}
}

func (v *Visualizer) formatNodes(w io.Writer) {
	prnDist := func(m *Multinomial) string {
		var buf bytes.Buffer
		for k, v := range m.Hist {
			fmt.Fprintf(&buf, "%s:%f ", k, v/m.Sum)
		}
		return buf.String()
	}

	prnChans := func(outDist []*Multinomial) string {
		var buf bytes.Buffer
		for _, m := range outDist {
			fmt.Fprintf(&buf, "%s\\n", prnDist(m))
		}
		return buf.String()
	}

	for state, channels := range v.Σγo {
		fmt.Fprintf(w, "%05d [shape=ellipse,label=\"%s\"];\n",
			state, prnChans(channels))
	}
}

func (v *Visualizer) formatEdges(w io.Writer, threshold float64) {
	for i := 0; i < len(v.Σξ); i++ {
		for j := 0; j < len(v.Σξ[i]); j++ {
			if p := v.Σξ[i][j]; p > 0 {
				fmt.Fprintf(w, "%05d -> %05d [label=\"%f\",weight=%d];\n",
					i, j, float64(p)/float64(v.Σγ[i]), int(p-threshold))
			}
		}
	}
}

// thresholdEdgeWeight sorts all weights order and finds top N weights
// that accumulate to frac of total weights; it then returns the
// weight of rank N.
func (v *Visualizer) thresholdEdgeWeight(frac float64) float64 {
	if frac < 0 || frac > 1 {
		panic(fmt.Sprintf("frac (%d) out of range [0,1]", frac))
	}

	w := make([]float64, 0)
	sum := 0.0
	for _, row := range v.Σξ {
		for _, v := range row {
			w = append(w, v)
			sum += v
		}
	}

	sort.Float64s(w)

	partial := 0.0
	for i := len(w) - 1; i >= 0; i-- {
		partial += w[i]
		if partial >= frac*sum {
			return w[i]
		}
	}
	return 0 // Display all edges.
}
