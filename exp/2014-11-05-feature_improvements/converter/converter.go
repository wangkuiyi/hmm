package main

import (
       "github.com/wangkuiyi/hmm/exp/linkedin_employee_data/converter.v3"
)

func main() {
     converter.Run(new(converter.PlainFeatureGenerator))
}