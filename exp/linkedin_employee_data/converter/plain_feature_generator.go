package converter

type PlainFeatureGenerator struct{}

func (*PlainFeatureGenerator) NumChannels() int {
	return 5
}

func (*PlainFeatureGenerator) Feature(r *Record, ch int) string {
	base := ""
	switch ch {
	case 0:
		base = r.CompanyOrSchool
	case 1:
		base = r.PosOrEduRank
	case 2:
		base = r.TitleOrDegree
	case 3:
		base = r.SeniorityOrDegreeRank
	case 4:
		base = r.FunctionOrField
	}

	// Some field values are negatives, which means nothing.
	if base[0] == '-' {
		return ""
	}

	return prefix(r, ch) + base
}

func prefix(r *Record, ch int) string {
	codebook := [][]string{
		{"company", "school"},
		{"position", "rank"},
		{"title", "degree"},
		{"seniority", "degreeRank"},
		{"function", "field"}}

	ret := ""
	if r.IsJob {
		ret += codebook[ch][0]
	}
	if r.IsEdu {
		ret += codebook[ch][1]
	}
	return ret
}
