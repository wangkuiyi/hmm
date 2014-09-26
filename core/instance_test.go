package core

import (
	"reflect"
	"testing"
)

func TestBuildInstanceIndex(t *testing.T) {
	index := buildInstanceIndex(kDachengPeriods)
	exp := 0
	for _, l := range kDachengPeriods {
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

func TestNewInstance(t *testing.T) {
	dacheng := NewInstance(kDachengObs, kDachengPeriods)
	if !reflect.DeepEqual(dacheng.Obs, kDachengObs) {
		t.Errorf("Expecting %v, got %v", kDachengObs, dacheng.Obs)
	}
	if !reflect.DeepEqual(dacheng.Periods, kDachengPeriods) {
		t.Errorf("Expecting %v, got %v", kDachengPeriods, dacheng.Periods)
	}

	null := NewInstance([][]Observed{nil, nil}, []int{0, 0})
	if null != nil {
		t.Errorf("Expecting nil, got %v", null)
	}
}
