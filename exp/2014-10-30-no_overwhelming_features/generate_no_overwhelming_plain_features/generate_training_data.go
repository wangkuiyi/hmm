package main

import (
	gen "github.com/wangkuiyi/hmm/exp/corpus_generation"
)

func main() {
	gen.Run(new(NoOverwhelmingFeatures))
}

type NoOverwhelmingFeatures struct{}

func (*NoOverwhelmingFeatures) NumChannels() int {
	return 1
}

// The following feature black list comes from
// hmm/exp/linkedin_employee_data/statistics/README.md.
/*
// company: 145411
// Education rank: 1, 2, 3
// Degree: 56, 22
// Seniority: 4, 3, 5
// Degree rank: 2, 3
// function: 25, 8
// field: 42
*/
func (*NoOverwhelmingFeatures) Feature(rs []*gen.Record, y, ch int) []string {
	ret := make([]string, 0)

	// NOTE: The following code is abandoned, as new versions of input
	// does not have fields like CompnayOrSchool any more.
	/*
		for _, r := range rs {
			if v := r.CompanyOrSchool; r.IsJob &&
				v != "145411" && v != "-9" {
				ret = append(ret, "company"+v)
			} else if r.IsEdu && v != "-9" {
				ret = append(ret, "school"+v)
			}

			// Do not use PosOrEduRank, which is not confident.

			if v := r.TitleOrDegree; r.IsEdu &&
				v != "56" && v != "22" && v != "-9" {
				ret = append(ret, "degree"+v)
			} else if r.IsJob && v != "-9" {
				ret = append(ret, "title"+v)
			}

			if v := r.SeniorityOrDegreeRank; r.IsJob &&
				v != "4" && v != "3" && v != "5" && v != "-9" {
				ret = append(ret, "seniority"+v)
			} else if r.IsEdu && v != "2" && v != "3" && v != "-9" {
				ret = append(ret, "degreerank"+v)
			}

			if v := r.FunctionOrField; r.IsJob &&
				v != "25" && v != "8" && v != "-9" {
				ret = append(ret, "function"+v)
			} else if r.IsEdu && (v != "42") && v != "-9" {
				ret = append(ret, "field"+v)
			}
		}
	*/
	return ret
}
