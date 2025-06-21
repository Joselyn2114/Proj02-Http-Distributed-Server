package main

import (
    "math/rand"
    "strconv"
    "time"

    "github.com/KateGF/Http-Server-Project-SO/core"
)

// piPartHandler fragmenta el cálculo de π según ?iter=n
func piPartHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    iterStr := req.Target.Query().Get("iter")
    iter, err := strconv.Atoi(iterStr)
    if err != nil || iter < 1 {
        return core.BadRequest().Text("invalid 'iter' parameter"), nil
    }

    rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
    inside := 0
    for i := 0; i < iter; i++ {
        x := rnd.Float64()
        y := rnd.Float64()
        if x*x+y*y <= 1 {
            inside++
        }
    }

    // Devolvemos {"inside": <count>}
    return core.Ok().JsonObj(map[string]int{"inside": inside}), nil
}
