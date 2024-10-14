package network

type Matrix [][]float64

type Network struct {
	Layers  []Matrix
	Weights []Matrix
	Biases  []Matrix
	Output  Matrix
	Rate    float64
	Errors  []float64
	Time    float64
	Locale  string
}

type Derivative struct {
	Delta      Matrix
	Adjustment Matrix
}
