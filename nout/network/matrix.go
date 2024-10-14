package network

import (
	"math/rand"
)

func CreateMatrix(rows, columns int) (matrix Matrix) {
	matrix = make(Matrix, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, columns)
	}

	return
}

func Columns(matrix Matrix) int {
	return len(matrix[0])
}

func RandomMatrix(rows, columns int) (matrix Matrix) {
	matrix = make(Matrix, rows)

	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, columns)
		for j := 0; j < columns; j++ {
			matrix[i][j] = rand.Float64()*2.0 - 1.0
		}
	}

	return
}

func Rows(matrix Matrix) int {
	return len(matrix)
}

func CreateNetwork(locale string, rate float64, input, output Matrix, hiddensNodes ...int) Network {
	input = append([][]float64{
		make([]float64, len(input[0])),
	}, input...)
	output = append([][]float64{
		make([]float64, len(output[0])),
	}, output...)

	inputMatrix := input
	layers := []Matrix{inputMatrix}

	for _, hiddenNodes := range hiddensNodes {
		layers = append(layers, CreateMatrix(len(input), hiddenNodes))
	}

	layers = append(layers, output)

	weightsNumber := len(layers) - 1
	var weights []Matrix
	var biases []Matrix

	for i := 0; i < weightsNumber; i++ {
		rows, columns := Columns(layers[i]), Columns(layers[i+1])

		weights = append(weights, RandomMatrix(rows, columns))
		biases = append(biases, RandomMatrix(Rows(layers[i]), columns))
	}

	return Network{
		Layers:  layers,
		Weights: weights,
		Biases:  biases,
		Output:  output,
		Rate:    rate,
		Locale:  locale,
	}
}

func DotProduct(matrix, matrix2 Matrix) Matrix {
	if Columns(matrix) != Rows(matrix2) {
		panic("Cannot make dot product between these two matrix.")
	}

	return ApplyFunctionWithIndex(
		CreateMatrix(Rows(matrix), Columns(matrix2)),
		func(i, j int, x float64) float64 {
			var sum float64

			for k := 0; k < Columns(matrix); k++ {
				sum += matrix[i][k] * matrix2[k][j]
			}

			return sum
		},
	)
}

func ErrorNotSameSize(matrix, matrix2 Matrix) {
	if Rows(matrix) != Rows(matrix2) && Columns(matrix) != Columns(matrix2) {
		panic("These two matrices must have the same dimension.")
	}
}

func ApplyFunctionWithIndex(matrix Matrix, fn func(i, j int, x float64) float64) Matrix {
	for i := 0; i < Rows(matrix); i++ {
		for j := 0; j < Columns(matrix); j++ {
			matrix[i][j] = fn(i, j, matrix[i][j])
		}
	}

	return matrix
}

func Sum(matrix, matrix2 Matrix) (resultMatrix Matrix) {
	ErrorNotSameSize(matrix, matrix2)

	resultMatrix = CreateMatrix(Rows(matrix), Columns(matrix)) // warning is fine (lint)

	return ApplyFunctionWithIndex(matrix, func(i, j int, x float64) float64 {
		return matrix[i][j] + matrix2[i][j]
	})
}

func ApplyFunction(matrix Matrix, fn func(x float64) float64) Matrix {
	return ApplyFunctionWithIndex(matrix, func(i, j int, x float64) float64 {
		return fn(x)
	})
}

func Differencen(matrix, matrix2 Matrix) (resultMatrix Matrix) {
	ErrorNotSameSize(matrix, matrix2)

	resultMatrix = CreateMatrix(Rows(matrix), Columns(matrix))

	return ApplyFunctionWithIndex(resultMatrix, func(i, j int, x float64) float64 {
		return matrix[i][j] - matrix2[i][j]
	})
}

func Multiplication(matrix, matrix2 Matrix) (resultMatrix Matrix) {
	ErrorNotSameSize(matrix, matrix2)

	resultMatrix = CreateMatrix(Rows(matrix), Columns(matrix)) // warning is fine (lint)

	return ApplyFunctionWithIndex(matrix, func(i, j int, x float64) float64 {
		return matrix[i][j] * matrix2[i][j]
	})
}

func Transpose(matrix Matrix) (resultMatrix Matrix) {
	resultMatrix = CreateMatrix(Columns(matrix), Rows(matrix))

	for i := 0; i < Rows(matrix); i++ {
		for j := 0; j < Columns(matrix); j++ {
			resultMatrix[j][i] = matrix[i][j]
		}
	}

	return resultMatrix
}

func ApplyRate(matrix Matrix, rate float64) Matrix {
	return ApplyFunction(matrix, func(x float64) float64 {
		return rate * x
	})
}
