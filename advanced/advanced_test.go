package advanced

import (
	"encoding/json"
	"github.com/KateGF/Http-Server-Project-SO/core"
	"net/url"
	"testing"
	"time"
)

func makeReq(path string) *core.HttpRequest {
	u, _ := url.Parse(path)
	return core.NewHttpRequest("GET", u, map[string]string{}, "")
}

func TestRandomHandler_Success(t *testing.T) {
	req := makeReq("/random?count=5&min=1&max=3")
	res, _ := RandomHandler(req)
	if res.StatusCode != 200 {
		t.Fatalf("want 200; got %d", res.StatusCode)
	}
	var body struct {
		Numbers []int `json:"numbers"`
	}
	if err := json.Unmarshal([]byte(res.Body), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(body.Numbers) != 5 {
		t.Errorf("want 5 nums; got %d", len(body.Numbers))
	}
	for _, v := range body.Numbers {
		if v < 1 || v > 3 {
			t.Errorf("value %d out of range", v)
		}
	}
}

func TestRandomHandler_Errors(t *testing.T) {
	cases := []struct {
		query, wantBody string
	}{
		{"", "count must be a positive integer"},
		{"count=0&min=1&max=3", "count must be a positive integer"},
		{"count=a&min=1&max=3", "count must be a positive integer"},
		{"count=3&max=3", "min must be a number"},
		{"count=3&min=a&max=3", "min must be a number"},
		{"count=3&min=1", "max must be a number"},
		{"count=3&min=5&max=2", "max must be >= min"},
	}

	for _, tc := range cases {
		req := makeReq("/random?" + tc.query)
		res, _ := RandomHandler(req)
		if res.StatusCode != 400 {
			t.Errorf("random?%s: want status 400; got %d", tc.query, res.StatusCode)
		}
		if res.Body != tc.wantBody {
			t.Errorf("random?%s: want body %q; got %q", tc.query, tc.wantBody, res.Body)
		}
	}
}

func TestTimestampHandler(t *testing.T) {
	req := makeReq("/timestamp")
	res, _ := TimestampHandler(req)
	if res.StatusCode != 200 {
		t.Fatalf("want 200; got %d", res.StatusCode)
	}
	var body struct {
		Timestamp string `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(res.Body), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Timestamp == "" {
		t.Error("empty timestamp")
	}
}

func TestSimulateHandler(t *testing.T) {
	// seconds=0 para no esperar realmente
	req := makeReq("/simulate?seconds=0&task=demo")
	res, _ := SimulateHandler(req)
	if res.StatusCode != 200 {
		t.Fatalf("want 200; got %d", res.StatusCode)
	}
	var body struct {
		Task string `json:"task"`
		Done bool   `json:"done"`
	}
	if err := json.Unmarshal([]byte(res.Body), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Task != "demo" || !body.Done {
		t.Errorf("unexpected body %+v", body)
	}
}

func TestSimulateHandler_Errors(t *testing.T) {
	cases := []struct {
		query, wantBody string
	}{
		{"", "seconds is required"},
		{"seconds=a&task=foo", "seconds must be a number"},
		{"seconds=-1&task=foo", "seconds must be >= 0"},
		{"seconds=0", "task is required"},
	}

	for _, tc := range cases {
		req := makeReq("/simulate?" + tc.query)
		res, _ := SimulateHandler(req)
		if res.StatusCode != 400 {
			t.Errorf("simulate?%s: want status 400; got %d", tc.query, res.StatusCode)
		}
		if res.Body != tc.wantBody {
			t.Errorf("simulate?%s: want body %q; got %q", tc.query, tc.wantBody, res.Body)
		}
	}
}

func TestSleepHandler_Success(t *testing.T) {
	req := makeReq("/sleep?seconds=0")
	res, _ := SleepHandler(req)
	if res.StatusCode != 200 {
		t.Fatalf("want 200; got %d", res.StatusCode)
	}
	expected := "slept 0 seconds"
	if res.Body != expected {
		t.Errorf("want body %q; got %q", expected, res.Body)
	}

	req = makeReq("/sleep?seconds=2")
	// Aquí podría tardar 2s, pero en test puedes usar seconds=0 ó mock Tiempo.
	res, _ = SleepHandler(req)
	if res.Body != "slept 2 seconds" {
		t.Errorf("want body %q; got %q", "slept 2 seconds", res.Body)
	}
}

func TestSleepHandler_Errors(t *testing.T) {
	cases := []struct{ query, want string }{
		{"", "seconds is required"},
		{"seconds=abc", "seconds must be a number"},
		{"seconds=-1", "seconds must be >= 0"},
	}
	for _, tc := range cases {
		req := makeReq("/sleep?" + tc.query)
		res, _ := SleepHandler(req)
		if res.StatusCode != 400 {
			t.Errorf("sleep?%s: want status 400; got %d", tc.query, res.StatusCode)
		}
		if res.Body != tc.want {
			t.Errorf("sleep?%s: want body %q; got %q", tc.query, tc.want, res.Body)
		}
	}
}

func TestLoadTestHandler_Success(t *testing.T) {
	// Con sleep=0 no tardará, así que es rápido
	req := makeReq("/loadtest?tasks=5&sleep=0")
	res, _ := LoadTestHandler(req)
	if res.StatusCode != 200 {
		t.Fatalf("want 200; got %d", res.StatusCode)
	}
	var body struct {
		Tasks      int   `json:"tasks"`
		Sleep      int   `json:"sleep"`
		DurationMS int64 `json:"duration_ms"`
	}
	if err := json.Unmarshal([]byte(res.Body), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if body.Tasks != 5 {
		t.Errorf("want tasks=5; got %d", body.Tasks)
	}
	if body.Sleep != 0 {
		t.Errorf("want sleep=0; got %d", body.Sleep)
	}
	if body.DurationMS < 0 {
		t.Errorf("want non-negative duration; got %d", body.DurationMS)
	}
}

func TestLoadTestHandler_Errors(t *testing.T) {
	cases := []struct {
		query, wantBody string
	}{
		{"", "tasks is required"},
		{"tasks=a&sleep=0", "tasks must be a number"},
		{"tasks=0&sleep=0", "tasks must be >= 1"},
		{"tasks=3", "sleep is required"},
		{"tasks=3&sleep=a", "sleep must be a number"},
		{"tasks=3&sleep=-1", "sleep must be >= 0"},
	}

	for _, tc := range cases {
		req := makeReq("/loadtest?" + tc.query)
		res, _ := LoadTestHandler(req)
		if res.StatusCode != 400 {
			t.Errorf("loadtest?%s: want status 400; got %d", tc.query, res.StatusCode)
		}
		if res.Body != tc.wantBody {
			t.Errorf("loadtest?%s: want body %q; got %q", tc.query, tc.wantBody, res.Body)
		}
	}
}

func TestStatusHandler(t *testing.T) {
	req := makeReq("/status")
	for i := 0; i < 3; i++ {
		// Simulamos una conexión
		CountConn()

		res, _ := StatusHandler(req)
		if res.StatusCode != 200 {
			t.Fatalf("status: want 200; got %d", res.StatusCode)
		}
		var body struct {
			Uptime     float64 `json:"uptime_s"`
			TotalConns int64   `json:"total_connections"`
			PID        int     `json:"pid"`
			Goroutines int     `json:"goroutines"`
		}
		if err := json.Unmarshal([]byte(res.Body), &body); err != nil {
			t.Fatalf("status JSON: %v", err)
		}

		// Ahora esperamos que TotalConns == i+1
		wantConns := int64(i + 1)
		if body.TotalConns != wantConns {
			t.Errorf("after %d CountConn(), got TotalConns=%d; want %d", i+1, body.TotalConns, wantConns)
		}
		if body.PID < 1 {
			t.Errorf("invalid PID %d", body.PID)
		}
		if body.Goroutines < 1 {
			t.Errorf("invalid Goroutines %d", body.Goroutines)
		}

		// esperamos uptime creciente
		if body.Uptime <= 0 {
			t.Errorf("uptime must be >0, got %f", body.Uptime)
		}

		// pequeña pausa
		time.Sleep(10 * time.Millisecond)
	}
}

func TestHelpHandler(t *testing.T) {
	req := makeReq("/help")
	res, _ := HelpHandler(req)
	if res.StatusCode != 200 {
		t.Fatalf("help: want 200; got %d", res.StatusCode)
	}
	var body struct {
		Commands []string `json:"commands"`
	}
	if err := json.Unmarshal([]byte(res.Body), &body); err != nil {
		t.Fatalf("help JSON: %v", err)
	}
	// nos basta que contenga al menos uno de los comandos conocidos
	found := false
	for _, cmd := range body.Commands {
		if cmd == "GET  /fibonacci?num=" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("help missing fibonacci command: %+v", body.Commands)
	}
}
