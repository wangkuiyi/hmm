package main

import (
	"github.com/wangkuiyi/hmm/exp/2014-10-27-linkedin-employees-plain-features/generate_training_data/generator"
)

func main() {
	generator.Run(new(generator.PlainFeatureGenerator))
}
