package core

var (
	kDachengObs = [][]Observed{
		[]Observed{
			Observed{"vice": 1, "president": 1},
			Observed{"applied": 1, "predictive": 1}},
		[]Observed{
			Observed{"head": 1, "solution": 1, "strategy": 1},
			Observed{"opera": 1, "solutions": 1}},
		[]Observed{
			Observed{"founder": 1},
			Observed{"helping": 1}},
		[]Observed{
			Observed{"senior": 1, "manager": 1},
			Observed{"linkedin": 1}}}
	kDachengPeriods = []int{8, 0, 1, 1}

	kGuanObs = [][]Observed{
		[]Observed{
			Observed{"cofounder": 1},
			Observed{"scissorsfly": 1}},
		[]Observed{
			Observed{"sr": 1, "associate": 1, "business": 1, "analyst": 1},
			Observed{"linkedin": 1}}}
	kGuanPeriods = []int{1, 1}

	kYiObs = [][]Observed{
		[]Observed{
			Observed{"software": 1, "engineer": 1},
			Observed{"google": 1}},
		[]Observed{
			Observed{"engineering": 1, "director": 1},
			Observed{"tencent": 1, "ads": 1}},
		[]Observed{
			Observed{"data": 1, "scientist": 1},
			Observed{"linkedin": 1}}}
	kYiPeriods = []int{3, 4, 0}

	kN = 2
	kC = 2
)

type mockRng struct {
	History []int
}

func (rng *mockRng) Intn(n int) int {
	if len(rng.History) == 0 {
		rng.History = make([]int, 1, 100)
		rng.History[0] = 0
		return 0
	}
	p := rng.History[len(rng.History)-1]
	if p+1 >= n {
		p = 0
	} else {
		p = p + 1
	}
	rng.History = append(rng.History, p)
	return p
}
