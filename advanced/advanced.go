package advanced

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

var (
	startTime  = time.Now()
	totalConns int64
)

// Invocado desde core.HandleWithError o en Middleware
func CountConn() {
	atomic.AddInt64(&totalConns, 1)
}

// SimulateHandler simula una tarea cuyo procesamiento toma 'seconds' segundos.
// URL: /simulate?seconds=s&task=name
func SimulateHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	// Extraer query params
	q := req.Target.Query()
	secStr := q.Get("seconds")
	if secStr == "" {
		return core.BadRequest().Text("seconds is required"), nil
	}
	seconds, err := strconv.Atoi(secStr)
	if err != nil {
		return core.BadRequest().Text("seconds must be a number"), nil
	}
	if seconds < 0 {
		return core.BadRequest().Text("seconds must be >= 0"), nil
	}

	task := q.Get("task")
	if task == "" {
		return core.BadRequest().Text("task is required"), nil
	}

	// Simular la tarea
	time.Sleep(time.Duration(seconds) * time.Second)

	// Construir respuesta JSON
	resp := struct {
		Task string `json:"task"`
		Done bool   `json:"done"`
	}{task, true}

	return core.Ok().JsonObj(resp), nil
}

// SleepHandler simula un retardo sin otra lógica.
// URL: /sleep?seconds=s
func SleepHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	q := req.Target.Query()
	secStr := q.Get("seconds")
	if secStr == "" {
		return core.BadRequest().Text("seconds is required"), nil
	}
	seconds, err := strconv.Atoi(secStr)
	if err != nil {
		return core.BadRequest().Text("seconds must be a number"), nil
	}
	if seconds < 0 {
		return core.BadRequest().Text("seconds must be >= 0"), nil
	}

	// Sleep real (o cero si seconds == 0)
	time.Sleep(time.Duration(seconds) * time.Second)

	// Respuesta simple en texto plano
	return core.Ok().Text(fmt.Sprintf("slept %d seconds", seconds)), nil
}

// LoadTestHandler simula 'tasks' goroutines durmiendo 'sleep' segundos cada una.
// URL: /loadtest?tasks=n&sleep=x
func LoadTestHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	q := req.Target.Query()
	tasksStr := q.Get("tasks")
	if tasksStr == "" {
		return core.BadRequest().Text("tasks is required"), nil
	}
	n, err := strconv.Atoi(tasksStr)
	if err != nil {
		return core.BadRequest().Text("tasks must be a number"), nil
	}
	if n < 1 {
		return core.BadRequest().Text("tasks must be >= 1"), nil
	}

	sleepStr := q.Get("sleep")
	if sleepStr == "" {
		return core.BadRequest().Text("sleep is required"), nil
	}
	x, err := strconv.Atoi(sleepStr)
	if err != nil {
		return core.BadRequest().Text("sleep must be a number"), nil
	}
	if x < 0 {
		return core.BadRequest().Text("sleep must be >= 0"), nil
	}

	// Lanzar n goroutines y medir tiempo
	var wg sync.WaitGroup
	wg.Add(n)
	start := time.Now()
	for i := 0; i < n; i++ {
		go func() {
			time.Sleep(time.Duration(x) * time.Second)
			wg.Done()
		}()
	}
	wg.Wait()
	durationMs := time.Since(start).Milliseconds()

	// Construir JSON de salida
	resp := struct {
		Tasks      int   `json:"tasks"`
		Sleep      int   `json:"sleep"`
		DurationMS int64 `json:"duration_ms"`
	}{n, x, durationMs}

	return core.Ok().JsonObj(resp), nil
}

// StatusHandler
func StatusHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	uptime := time.Since(startTime).Seconds()
	resp := struct {
		Uptime     float64 `json:"uptime_s"`
		TotalConns int64   `json:"total_connections"`
		PID        int     `json:"pid"`
		Goroutines int     `json:"goroutines"`
	}{
		uptime,
		atomic.LoadInt64(&totalConns),
		os.Getpid(),
		runtime.NumGoroutine(),
	}
	return core.Ok().JsonObj(resp), nil
}

// HelpHandler (/help)
func HelpHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	cmds := []string{
		"GET  /fibonacci?num=",
		"POST /createfile?name=&content=&repeat=",
		"DELETE /deletefile?name=",
		"GET  /reverse?text=",
		"GET  /toupper?text=",
		"GET  /hash?text=",
		"GET  /random?count=&min=&max=",
		"GET  /timestamp",
		"GET  /simulate?seconds=&task=",
		"GET  /sleep?seconds=",
		"GET  /loadtest?tasks=&sleep=",
		"GET  /status",
		"GET  /help",
	}
	return core.Ok().JsonObj(struct {
		Commands []string `json:"commands"`
	}{cmds}), nil
}

// RandomHandler: /random?count=n&min=a&max=b
func RandomHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	q := req.Target.Query()

	// count
	cnt, err := strconv.Atoi(q.Get("count"))
	if err != nil || cnt < 1 {
		return core.BadRequest().Text("count must be a positive integer"), nil
	}
	// min
	min, err := strconv.Atoi(q.Get("min"))
	if err != nil {
		return core.BadRequest().Text("min must be a number"), nil
	}
	// max
	max, err := strconv.Atoi(q.Get("max"))
	if err != nil {
		return core.BadRequest().Text("max must be a number"), nil
	}
	if max < min {
		return core.BadRequest().Text("max must be >= min"), nil
	}

	// genera números
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		nums[i] = rand.Intn(max-min+1) + min
	}

	// devuelve JSON {"numbers":[...]}
	return core.Ok().JsonObj(struct {
		Numbers []int `json:"numbers"`
	}{nums}), nil
}

// TimestampHandler: /timestamp
func TimestampHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	return core.Ok().JsonObj(struct {
		Timestamp string `json:"timestamp"`
	}{now}), nil
}
