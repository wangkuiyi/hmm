package core

import (
	"reflect"
	"testing"
)

var (
	dacheng = NewInstance(
		[][]Observed{
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
				Observed{"linkedin": 1}}},
		[]int{8, 0, 1, 1})

	yiwang = NewInstance(
		[][]Observed{
			[]Observed{
				Observed{"software": 1, "engineer": 1},
				Observed{"google": 1}},
			[]Observed{
				Observed{"engineering": 1, "director": 1},
				Observed{"tencent": 1, "ads": 1}},
			[]Observed{
				Observed{"data": 1, "scientist": 1},
				Observed{"linkedin": 1}}},
		[]int{3, 4, 0})

	gwang = NewInstance(
		[][]Observed{
			[]Observed{
				Observed{"cofounder": 1},
				Observed{"scissorsfly": 1}},
			[]Observed{
				Observed{"sr": 1, "associate": 1, "business": 1, "analyst": 1},
				Observed{"linkedin": 1}}},
		[]int{1, 1})
)

func TestBuildInstanceIndex(t *testing.T) {
	index := buildInstanceIndex(dacheng.Periods)
	exp := 0
	for _, l := range dacheng.Periods {
		exp += l
	}
	if len(index) != exp {
		t.Errorf("Expecting %d, got %d", exp, len(index))
	}
	truth := []int{0, 0, 0, 0, 0, 0, 0, 0, 2, 3}
	if !reflect.DeepEqual(index, truth) {
		t.Errorf("Expecting %d, got %d", truth, index)
	}
}

// func TestTrain(t *testing.T) {
// 	corpus := []*Instance{}

// 	C := EstimateC(corpus)
// 	K := 2
// 	Iter := 100
// 	model := Train(corpus, K, C, Iter)
// 	fmt.Println(model)
// }
