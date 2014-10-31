package main

import (
	"flag"
	"github.com/wangkuiyi/hmm/core"
	"log"
)

func main() {
	flagModel := flag.String("model", "", "Model file in JSON")
	flagFigure := flag.String("figure", "", "Figure file in PDF")
	flagEdge := flag.Float64("edge", 0.0, "Edge prunning threshold")
	flagNode := flag.Float64("node", 0.0, "Node prunning threshold")
	flag.Parse()

	m, e := core.LoadModel(*flagModel)
	if e != nil {
		log.Fatalf("Canont load model from %s: %v", *flagModel, e)
	}

	e = core.NewVisualizer(m).Draw(*flagFigure, *flagEdge, *flagNode)
	if e != nil {
		log.Fatalf("Cannot visualzie model: %v", e)
	}
}
