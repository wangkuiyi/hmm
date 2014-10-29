package main

import (
	"github.com/wangkuiyi/hmm/exp/linkedin_employee_data/converter"
)

func main() {
	converter.Run(new(converter.PlainFeatureGenerator))
}
