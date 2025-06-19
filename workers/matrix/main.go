package main

import (
	"os"
	"strconv"
	"sync/atomic"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func main() {
	var currentLoad int64
	var doneCount int64

	srv := core.NewHttpServer()

	srv.Get("/ping", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().JsonObj(map[string]any{
			"status": "ok",
			"load":   atomic.LoadInt64(&currentLoad),
			"done":   atomic.LoadInt64(&doneCount),
		}), nil
	})

	srv.Post("/matrix", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		atomic.AddInt64(&currentLoad, 1)
		defer atomic.AddInt64(&currentLoad, -1)
		response, err := MatrixHandler(req)
		if err == nil && response.StatusCode == 200 {
			atomic.AddInt64(&doneCount, 1)
		}
		return response, err
	})

	port := mustAtoi(os.Getenv("PORT"))
	srv.Start(port)
}
