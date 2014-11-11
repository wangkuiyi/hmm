package main

import (
       "github.com/chqsark/hmm/exp/linkedin_employee_data/converter.v3"
)

func main() {
     converter.Run(new(converter.PlainFeatureGenerator))
}
