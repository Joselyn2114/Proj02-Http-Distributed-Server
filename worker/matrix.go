package main

import (
    "encoding/json"

    "github.com/KateGF/Http-Server-Project-SO/core"
)

// payload para un subbloque de multiplicación
type matrixBlock struct {
    A [][]float64 `json:"a"`
    B [][]float64 `json:"b"`
}

func matrixPartHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    var blk matrixBlock
    if err := json.Unmarshal([]byte(req.Body), &blk); err != nil {
        return core.BadRequest().Text("invalid JSON"), nil
    }
    // multiplicación de matrices A×B
    rA, cA := len(blk.A), len(blk.A[0])
    _, cB := len(blk.B), len(blk.B[0])
    C := make([][]float64, rA)
    for i := 0; i < rA; i++ {
        C[i] = make([]float64, cB)
        for j := 0; j < cB; j++ {
            sum := 0.0
            for k := 0; k < cA; k++ {
                sum += blk.A[i][k] * blk.B[k][j]
            }
            C[i][j] = sum
        }
    }
    return core.Ok().JsonObj(C), nil
}
