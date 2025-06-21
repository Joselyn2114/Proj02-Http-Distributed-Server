package main

import (
    "log/slog"

    "github.com/KateGF/Http-Server-Project-SO/core"
    "github.com/KateGF/Http-Server-Project-SO/handlers"
    "github.com/KateGF/Http-Server-Project-SO/service"
    "github.com/KateGF/Http-Server-Project-SO/advanced"
)

// ping responde “pong”
func pingHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    return core.Ok().Text("pong"), nil
}

func main() {
    server := core.NewHttpServer()

    // Rutas originales (tal como en tu main.go)
    server.Get("/fibonacci", service.FibonacciHandler)
    server.Post("/createfile", service.CreateFileHandler)
    server.Get("/createfile", service.CreateFileHandler)
    server.Delete("/deletefile", service.DeleteFileHandler)
    server.Get("/deletefile", service.DeleteFileHandler)

    server.Get("/reverse", handlers.ReverseHandler)
    server.Get("/toupper", handlers.ToUpperHandler)
    server.Get("/hash", handlers.HashHandler)
    server.Get("/", handlers.RootHandler)

    server.Get("/random", advanced.RandomHandler)
    server.Get("/timestamp", advanced.TimestampHandler)
    server.Get("/simulate", advanced.SimulateHandler)
    server.Get("/sleep", advanced.SleepHandler)
    server.Get("/loadtest", advanced.LoadTestHandler)
    server.Get("/status", advanced.StatusHandler)
    server.Get("/help", advanced.HelpHandler)

    // --- Nuevos endpoints para procesamiento distribuido ---
    server.Get("/ping", pingHandler)
    server.Get("/pi/part", piPartHandler)            // definido más abajo
    server.Post("/matrix/part", matrixPartHandler)   // definido más abajo

    slog.Info("Worker arrancado en :8080")
    if err := server.Start(8080); err != nil {
        slog.Error("Worker error", "err", err)
    }
}
