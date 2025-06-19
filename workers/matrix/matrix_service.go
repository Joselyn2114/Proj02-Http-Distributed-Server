package main

import (
	"encoding/json"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

// Estructura que representa dos matrices para operaciones.
type Matrices struct {
	A Matrix `json:"A"`
	B Matrix `json:"B"`
}

// ReadMatrices lee el cuerpo de la solicitud HTTP y deserializa las matrices.
func ReadMatrices(req *core.HttpRequest) (Matrices, error) {
	var matrices Matrices

	err := json.Unmarshal([]byte(req.Body), &matrices)
	if err != nil {
		return Matrices{}, err
	}

	return matrices, nil
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
