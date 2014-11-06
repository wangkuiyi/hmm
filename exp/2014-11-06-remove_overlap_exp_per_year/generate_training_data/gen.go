package main

import (
	gen "github.com/wangkuiyi/hmm/exp/corpus_generation"
	"log"
	"strconv"
)

type V4DataGenerator struct{}

func (*V4DataGenerator) NumChannels() int {
	return 1
}

func (*V4DataGenerator) Feature(rs []*gen.Record, year, ch int) []string {
	ret := make([]string, 0, 4)

	topSeniority := -1
	mostSenior := -1
	for i, r := range rs {
		if s, e := strconv.Atoi(r.Seniority); e != nil {
			log.Fatalf("Cannot parse seniority of record %v", r)
		} else if s > topSeniority {
			topSeniority = s
			mostSenior = i
		}
	}

	if mostSenior < 0 {
		log.Fatalf("Cannot find the most senior experience for %s in year %s",
			rs[0].Member, year)
	}

	r := rs[mostSenior]
	if r.Company != "" {
		ret = append(ret, "company="+r.Company)
	}
	if r.Title != "" {
		ret = append(ret, "title="+r.Title)
	}
	if r.Seniority != "" {
		ret = append(ret, "seniority="+r.Seniority)
	}
	if r.Function != "" {
		ret = append(ret, "function="+r.Function)
	}

	if f := r.School + r.Degree + r.DegreeRank + r.Field; f != "" {
		log.Fatalf("Record %v has educational properties.", r)
	}

	return ret
}

func main() {
	gen.Run(new(V4DataGenerator))
}
