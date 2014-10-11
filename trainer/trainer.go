package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wangkuiyi/hmm/core"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	flagAddr := flag.String("addr", ":6060", "Listening address")
	flagCorpus := flag.String("corpus", "", "Corpus file in JSON format")
	flagStates := flag.Int("states", 2, "Number of hidden states")
	flagIter := flag.Int("iter", 20, "Number of EM iterations")
	flagModel := flag.String("model", "", "Model file in JSON format")
	flagLL := flag.String("logl", "", "Log-likelihood file")
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe(*flagAddr, nil))
	}()

	var corpus []*core.Instance
	if f, e := os.Open(*flagCorpus); e != nil {
		log.Fatalf("Cannot open %s: %v", *flagCorpus, e)
	} else {
		defer f.Close()
		if e := json.NewDecoder(f).Decode(&corpus); e != nil {
			log.Fatalf("Failed decode corpus: %v", e)
		}
		// Infer unexported fileds of Instance.
		for i, _ := range corpus {
			corpus[i].Index()
		}
	}

	C := core.EstimateC(corpus)
	baseline := core.Init(*flagStates, C, corpus, rand.New(rand.NewSource(99)))
	model, ll := core.Train(corpus, *flagStates, C, *flagIter, baseline)

	core.SaveModel(model, *flagModel)

	f := core.CreateFileOrStdout(*flagLL)
	if f != os.Stdout {
		defer f.Close()
	}
	for _, l := range ll {
		fmt.Fprintln(f, l)
	}
}
