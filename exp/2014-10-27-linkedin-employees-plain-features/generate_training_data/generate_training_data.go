package main

import (
	corpus "github.com/wangkuiyi/hmm/exp/corpus_generation"
)

func main() {
	corpus.Run(new(corpus.PlainFeatureGenerator))
}
