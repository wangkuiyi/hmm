package main

import (
	"flag"
	"github.com/wangkuiyi/hmm/core"
	"log"
)

func main() {
	flagModel := flag.String("model", "", "Model file in JSON")
	flagFigure := flag.String("figure", "", "Figure file in PDF")
	flag.Parse()

	m, e := core.LoadModel(*flagModel)
	if e != nil {
		log.Fatalf("Canont load model from %s: %v", *flagModel, e)
	}

	if e := core.NewVisualizer(m).Draw(*flagFigure); e != nil {
		log.Fatalf("Cannot visualzie model: %v", e)
	}
}
