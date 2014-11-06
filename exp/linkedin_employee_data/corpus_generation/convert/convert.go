package main

import (
	cvt "github.com/wangkuiyi/hmm/exp/linkedin_employee_data/converter.v3"
)

func main() {
	cvt.Run(new(cvt.PlainFeatureGenerator))
}
