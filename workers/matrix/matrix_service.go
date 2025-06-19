package main

import (
	"encoding/json"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

// Estructura que representa dos matrices para operaciones.
type Matrices struct {
	A Matrix
	B Matrix
}

// ReadMatrices lee el cuerpo de la solicitud HTTP, deserializa las matrices y crea objetos Matrix.
func ReadMatrices(req *core.HttpRequest) (Matrices, error) {
	var data struct {
		A [][]float64 `json:"A"`
		B [][]float64 `json:"B"`
	}
	if err := json.Unmarshal([]byte(req.Body), &data); err != nil {
		return Matrices{}, err
	}

	A, err := NewMatrix(data.A)
	if err != nil {
		return Matrices{}, err
	}

	B, err := NewMatrix(data.B)
	if err != nil {
		return Matrices{}, err
	}

	return Matrices{A, B}, nil
}

// MatrixHandler maneja la solicitud HTTP para multiplicar dos matrices.
func MatrixHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	matrices, err := ReadMatrices(req)
	if err != nil {
		return core.BadRequest().Text(err.Error()), nil
	}

	matrix, err := matrices.A.Multiply(matrices.B)
	if err != nil {
		return core.BadRequest().Text(err.Error()), nil
	}

	jsonData, err := matrix.ToJson()
	if err != nil {
		return core.NewHttpResponse(500, "Internal Server Error", err.Error()), nil
	}

	return core.Ok().Json(jsonData), nil
}
