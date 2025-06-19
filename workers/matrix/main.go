package main

import (
	"os"
	"strconv"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func main() {
	srv := core.NewHttpServer()

	srv.Get("/ping", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		return core.Ok().JsonObj(map[string]any{
			"status": "ok", "load": 0, "done": 0,
		}), nil
	})

	srv.Get("/matrix", MatrixHandler)

	port := mustAtoi(os.Getenv("PORT"))
	srv.Start(port)
}
