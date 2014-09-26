package core

import (
	"log"
)

type Instance struct {
	Obs     [][]Observed // A sequence of observed.
	Periods []int        // Periods of above observed.
	index   []int        // Map time t to an observed.
}

type Observed map[string]int

func NewInstance(obs [][]Observed, periods []int) *Instance {
	if len(periods) != len(obs) {
		log.Panicf("len(period)=%d, len(obs)=%d", len(periods), len(obs))
	}

	ret := &Instance{
		Obs:     obs,
		Periods: periods,
		index:   buildInstanceIndex(periods)}

	if len(ret.index) == 0 {
		log.Printf("Instance ignored with periods=0; obs: %v", obs)
		return nil
	}
	return ret
}

func (i *Instance) Index() {
	i.index = buildInstanceIndex(i.Periods)
}

func buildInstanceIndex(periods []int) []int {
	var T int
	for _, l := range periods {
		T += l
	}
	if T <= 0 {
		return nil
	}

	ret := make([]int, 0, T)
	for i, l := range periods {
		for j := 0; j < l; j++ {
			ret = append(ret, i)
		}
	}
	return ret
}

func (i *Instance) T() int {
	return len(i.index)
}

func (i *Instance) O(t int) []Observed {
	return i.Obs[i.index[t]]
}
