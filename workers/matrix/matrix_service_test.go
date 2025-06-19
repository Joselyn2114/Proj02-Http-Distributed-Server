package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

func TestReadMatrices_Success(t *testing.T) {
	// Arrange
	body := `{"A":[[1,2],[3,4]],"B":[[5,6],[7,8]]}`
	req := &core.HttpRequest{Body: body}
	expectedA := [][]float64{{1, 2}, {3, 4}}
	expectedB := [][]float64{{5, 6}, {7, 8}}

	// Act
	mats, err := ReadMatrices(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(mats.A.data, expectedA) {
		t.Errorf("expected A %v, got %v", expectedA, mats.A.data)
	}
	if !reflect.DeepEqual(mats.B.data, expectedB) {
		t.Errorf("expected B %v, got %v", expectedB, mats.B.data)
	}
}

func TestReadMatrices_InvalidJSON(t *testing.T) {
	// Arrange
	req := &core.HttpRequest{Body: "not json"}

	// Act
	_, err := ReadMatrices(req)

	// Assert
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestReadMatrices_InvalidMatrix(t *testing.T) {
	// Arrange
	req := &core.HttpRequest{Body: `{"A":[],"B":[[1]]}`}

	// Act
	_, err := ReadMatrices(req)

	// Assert
	if err == nil {
		t.Error("expected error for empty matrix A, got nil")
	}
}

func TestMatrixHandler_Success(t *testing.T) {
	// Arrange
	body := `{"A":[[1,2],[3,4]],"B":[[5,6],[7,8]]}`
	req := &core.HttpRequest{Body: body}
	expected := [][]float64{{19, 22}, {43, 50}}

	// Act
	resp, err := MatrixHandler(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	var data [][]float64
	if err := json.Unmarshal([]byte(resp.Body), &data); err != nil {
		t.Fatalf("invalid JSON body: %v", err)
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}

func TestMatrixHandler_InvalidJSON(t *testing.T) {
	// Arrange
	req := &core.HttpRequest{Body: "not json"}

	// Act
	resp, err := MatrixHandler(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestMatrixHandler_DimMismatch(t *testing.T) {
	// Arrange
	body := `{"A":[[1,2]],"B":[[1]]}`
	req := &core.HttpRequest{Body: body}

	// Act
	resp, err := MatrixHandler(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}
}
