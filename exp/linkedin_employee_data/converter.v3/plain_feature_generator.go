package converter

type PlainFeatureGenerator struct{}

func (*PlainFeatureGenerator) NumChannels() int {
	return 1
}

func (*PlainFeatureGenerator) Feature(r *Record, year, ch int) []string {
	ret := make([]string, 0, kNumFields)

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

	if r.School != "" {
		ret = append(ret, "school="+r.School)
	}
	if r.Degree != "" {
		ret = append(ret, "degree="+r.Degree)
	}
	if r.DegreeRank != "" {
		ret = append(ret, "degreeRank="+r.DegreeRank)
	}
	if r.Field != "" {
		ret = append(ret, "field="+r.Field)
	}

	return ret
}
