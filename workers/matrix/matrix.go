package main

import (
	"encoding/json"
	"errors"
)

// Matrix representa una matriz 2D con métodos para multiplicación y serialización JSON.
type Matrix struct {
	data [][]float64
}

// NewMatrix crea una nueva instancia de Matrix a partir de los datos proporcionados.
func NewMatrix(data [][]float64) (Matrix, error) {
	if len(data) == 0 || len(data[0]) == 0 {
		return Matrix{}, errors.New("the matrix can't be empty")
	}

	return Matrix{data}, nil
}

// Multiply realiza la multiplicación de dos matrices y devuelve una nueva matriz resultante.
func (a Matrix) Multiply(b Matrix) (Matrix, error) {
	if len(a.data[0]) != len(b.data) {
		return Matrix{}, errors.New("the number of columns in the first matrix must be equal to the number of rows in the second matrix")
	}

	result := make([][]float64, len(a.data))

	for i := range result {
		result[i] = make([]float64, len(b.data[0]))

		for j := range result[i] {
			for k := range b.data {
				result[i][j] += a.data[i][k] * b.data[k][j]
			}
		}
	}

	return NewMatrix(result)
}

// ToJson serializa la matriz a una cadena JSON.
func (m Matrix) ToJson() (string, error) {
	jsonData, err := json.Marshal(m.data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
