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

	fmt.Fprintf(f, "digraph Model {\n")
	v.formatInits(f)
	v.formatNodes(f)
	v.formatEdges(f, v.thresholdEdgeWeight(1))
	fmt.Fprintf(f, "}\n")
	return nil
}

func pct(x float64) string {
	if int(x*1000) < 1 {
		return ""
	}
	return fmt.Sprintf("%2.1f%%", x*100.0)
}

func (v *Visualizer) formatInits(w io.Writer) {
	fmt.Fprintf(w, "start [shape=box];\n")
	for i, p := range v.S1 {
		if l := pct(p / v.S1Sum); len(l) > 0 {
			fmt.Fprintf(w, "start -> %05d [label=\"%s\",weight=%d];\n",
				i, l, int(p))
		}
	}
}

type WeightedString struct {
	key    string
	weight float64
}

type WeightedStringSlice []WeightedString

func (ws WeightedStringSlice) Len() int {
	return len(ws)
}

func (ws WeightedStringSlice) Less(i, j int) bool {
	if ws[i].weight == ws[j].weight {
		return ws[i].key < ws[j].key
	}
	return ws[i].weight > ws[j].weight
}

func (ws WeightedStringSlice) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

func (v *Visualizer) formatNodes(w io.Writer) {
	prnDist := func(m *Multinomial) string {
		s := make(WeightedStringSlice, 0, len(m.Hist))
		for k, v := range m.Hist {
			s = append(s, WeightedString{k, v})
		}
		sort.Sort(s)

		var buf bytes.Buffer
		for _, w := range s {
			if l := pct(w.weight / m.Sum); len(l) > 0 {
				fmt.Fprintf(&buf, "%s:%s ", w.key, l)
			}
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
			if l := pct(v.Σξ[i][j] / v.Σγ[i]); len(l) > 0 {
				fmt.Fprintf(w, "%05d -> %05d [label=\"%s\",weight=%d];\n",
					i, j, l, int(v.Σξ[i][j]/v.Σγ[i]-threshold))
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
	for s, row := range v.Σξ {
		for _, tr := range row {
			w = append(w, tr/v.Σγ[s])
			sum += tr / v.Σγ[s]
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
