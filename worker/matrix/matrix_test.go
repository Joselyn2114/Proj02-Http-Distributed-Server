package matrix

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewMatrix_Empty(t *testing.T) {
	// Arrange and Act
	_, err := NewMatrix([][]float64{})

	// Assert
	if err == nil {
		t.Error("expected error for empty matrix, got nil")
	}
}

func TestMultiply_Success(t *testing.T) {
	// Arrange
	a, _ := NewMatrix([][]float64{{1, 2, 3}, {4, 5, 6}})
	b, _ := NewMatrix([][]float64{{7, 8}, {9, 10}, {11, 12}})

	expected := [][]float64{{58, 64}, {139, 154}}

	// Act
	result, err := a.Multiply(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert
	if !reflect.DeepEqual(result.data, expected) {
		t.Errorf("expected %v, got %v", expected, result.data)
	}
}

func TestMultiply_Error(t *testing.T) {
	// Arrange
	a, _ := NewMatrix([][]float64{{1, 2}})
	b, _ := NewMatrix([][]float64{{1, 2}, {3, 4}, {5, 6}})

	// Act
	_, err := a.Multiply(b)

	// Assert
	if err == nil {
		t.Error("expected error for dimension mismatch, got nil")
	}
}

func TestToJson(t *testing.T) {
	// Arrange
	m, _ := NewMatrix([][]float64{{1.5, 2.5}, {3.5, 4.5}})

	expected := [][]float64{{1.5, 2.5}, {3.5, 4.5}}

	// Act
	jsonStr, err := m.ToJson()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var data [][]float64
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if !reflect.DeepEqual(data, expected) {
		t.Errorf("expected %v, got %v", expected, data)
	}
}
