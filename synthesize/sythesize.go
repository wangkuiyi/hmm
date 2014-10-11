package main

import (
	"encoding/json"
	"flag"
	"github.com/wangkuiyi/hmm/core"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main() {
	flagModel := flag.String("model", "", "Model file in JSON format")
	flagInstances := flag.Int("instances", 10, "# synthesized instances")
	flagLength := flag.Int("length", 10, "Length of each instance")
	flagCardi := flag.Int("cardi", 4, "Cardinality of multinomial outputs")
	flagCorpus := flag.String("corpus", "", "Synthetic corpus file")
	flag.Parse()

	var m *core.Model
	var e error
	if m, e = core.LoadModel(*flagModel); e != nil {
		log.Printf("Cannot load %s: %v. Use default model.", *flagModel, e)
		m = new(core.Model)
		if e := json.NewDecoder(
			strings.NewReader(defaultModel)).Decode(m); e != nil {
			log.Fatalf("Cannot decode default model")
		}
	}

	f := core.CreateFileOrStdout(*flagCorpus)
	if f != os.Stdout {
		defer f.Close()
	} else {
		log.Printf("Cannot create file %s. Use stdout", *flagCorpus)
	}

	if e := json.NewEncoder(f).Encode(
		m.Sample(*flagInstances, *flagLength, *flagCardi,
			rand.New(rand.NewSource(99)))); e != nil {
		log.Fatalf("Cannot JSON-encode corpus, %v", e)
	}
}

var (
	defaultModel = `{
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
