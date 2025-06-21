package service

import (
	"fmt"
	"github.com/KateGF/Http-Server-Project-SO/core"
	"net/url"
	"testing"
)

func TestFibonacci(t *testing.T) {
	var tests = []struct {
		num      int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
		{20, 6765},
		{92, 7540113804746346429},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Fibonacci(%d)", tt.num), func(t *testing.T) {
			act := Fibonacci(tt.num)

			if act != tt.expected {
				t.Errorf("Expected Fibonacci(%d) to be %d, not %d", tt.num, tt.expected, act)
			}
		})
	}
}

func TestFibonacciHandler(t *testing.T) {
	tests := []struct {
		num      string
		expected string
	}{
		{"", "num is required"},
		{"A", "num must be a number"},
		{"-1", "num must be between 0 and 92"},
		{"93", "num must be between 0 and 92"},
		{"0", "0"},
		{"1", "1"},
		{"10", "55"},
		{"92", "7540113804746346429"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("FibonacciHandler %d", i), func(t *testing.T) {
			// Arrange
			target, _ := url.Parse(fmt.Sprintf("/fibonacci?num=%s", tt.num))

			request := core.NewHttpRequest("GET", target, map[string]string{}, "")

			// Act
			response, err := FibonacciHandler(request)

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, %v", err)
			}

			if response.Body != tt.expected {
				t.Errorf("Expected body to be %s, not %s", tt.expected, response.Body)
			}
		})
	}
}
