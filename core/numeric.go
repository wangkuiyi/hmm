package core

func vector(x int) []float64 {
	return make([]float64, x)
}

func matrix(x, y int) [][]float64 {
	ret := make([][]float64, x)
	for i, _ := range ret {
		ret[i] = make([]float64, y)
	}
	return ret
}
