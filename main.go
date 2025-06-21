package main

import (
	"github.com/KateGF/Http-Server-Project-SO/advanced"
	"github.com/KateGF/Http-Server-Project-SO/core"
	"github.com/KateGF/Http-Server-Project-SO/handlers"
	"github.com/KateGF/Http-Server-Project-SO/service"
	"log/slog"
)

func main() {
	// Crea una nueva instancia del servidor HTTP.
	server := core.NewHttpServer()

	// Registra un manejador para la ruta GET "/fibonacci".
	server.Get("/fibonacci", service.FibonacciHandler)

	// Registra un manejador para POST "/createfile"
	server.Post("/createfile", service.CreateFileHandler)
	// También exponer "/createfile" por GET para pruebas manuales sin body.
	server.Get("/createfile", service.CreateFileHandler)

	// Registra un manejador para DELETE "/deletefile"
	server.Delete("/deletefile", service.DeleteFileHandler)
	// También exponer "/deletefile" por GET para pruebas manuales sin body.
	server.Get("/deletefile", service.DeleteFileHandler)

	// Endpoints de cadenas
	server.Get("/reverse", handlers.ReverseHandler)
	server.Get("/toupper", handlers.ToUpperHandler)
	server.Get("/hash", handlers.HashHandler)
	server.Get("/", handlers.RootHandler)

	// Endpoints avanzados
	server.Get("/random", advanced.RandomHandler)
	server.Get("/timestamp", advanced.TimestampHandler)
	server.Get("/simulate", advanced.SimulateHandler)
	server.Get("/sleep", advanced.SleepHandler)
	server.Get("/loadtest", advanced.LoadTestHandler)
	server.Get("/status", advanced.StatusHandler)
	server.Get("/help", advanced.HelpHandler)

	// Inicia el servidor en el puerto 8080.
	err := server.Start(8080)

	if err != nil {
		slog.Error("Error starting or running server", "error", err)
	}
}
