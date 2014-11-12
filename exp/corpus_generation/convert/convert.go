package main

import (
	cvt "github.com/wangkuiyi/hmm/exp/corpus_generation"
)

func main() {
	cvt.Run(new(cvt.PlainFeatureGenerator))
}
