// integration/advanced_integration_test.go
package integration

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "testing"
    "time"
)

// helper: envía raw HTTP request y devuelve únicamente el body (asume status 200 OK).
func sendRaw(t *testing.T, request string) string {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        t.Fatalf("Dial failed: %v", err)
    }
    defer conn.Close()

    fmt.Fprint(conn, request)
    reader := bufio.NewReader(conn)

    // Leer status line y validar 200
    status, _ := reader.ReadString('\n')
    if !strings.Contains(status, "200") {
        t.Fatalf("Expected 200 OK, got %q", status)
    }
    // Saltar headers
    for {
        line, _ := reader.ReadString('\n')
        if line == "\r\n" {
            break
        }
    }
    // Leer body completo (una línea JSON)
    body, _ := reader.ReadString('\n')
    return body
}

func TestRandomEndpoint(t *testing.T) {
    body := sendRaw(t, "GET /random?count=3&min=5&max=10 HTTP/1.0\r\nHost: test\r\n\r\n")
    // Ahora el servidor responde {"numbers":[...]}
    if !strings.Contains(body, `"numbers":[`) {
        t.Errorf("Random response missing numbers field: %q", body)
    }
}

func TestTimestampEndpoint(t *testing.T) {
    body := sendRaw(t, "GET /timestamp HTTP/1.0\r\nHost: test\r\n\r\n")
    if !strings.Contains(body, "T") || !strings.Contains(body, "Z") {
        t.Errorf("Timestamp format incorrect: %q", body)
    }
}

func TestSimulateAndSleepEndpoints(t *testing.T) {
    start := time.Now()
    _ = sendRaw(t, "GET /sleep?seconds=1 HTTP/1.0\r\nHost: test\r\n\r\n")
    if dur := time.Since(start); dur < time.Second {
        t.Errorf("Sleep less than 1s: %v", dur)
    }

    start = time.Now()
    _ = sendRaw(t, "GET /simulate?seconds=1&task=test HTTP/1.0\r\nHost: test\r\n\r\n")
    if dur := time.Since(start); dur < time.Second {
        t.Errorf("Simulate less than 1s: %v", dur)
    }
}

func TestLoadtestEndpoint(t *testing.T) {
    body := sendRaw(t, "GET /loadtest?tasks=5&sleep=0 HTTP/1.0\r\nHost: test\r\n\r\n")
    // Respuesta debe contener tasks y duration_ms
    if !strings.Contains(body, `"tasks":5`) || !strings.Contains(body, `"duration_ms":`) {
        t.Errorf("Loadtest response unexpected: %q", body)
    }
}

func TestStatusAndHelpEndpoints(t *testing.T) {
    help := sendRaw(t, "GET /help HTTP/1.0\r\nHost: test\r\n\r\n")
    if !strings.Contains(help, "/fibonacci") || !strings.Contains(help, "/help") {
        t.Errorf("Help response missing commands: %q", help)
    }

    status := sendRaw(t, "GET /status HTTP/1.0\r\nHost: test\r\n\r\n")
    // Ahora campos: uptime_s, total_connections, goroutines
    if !strings.Contains(status, `"uptime_s"`) ||
       !strings.Contains(status, `"total_connections"`) ||
       !strings.Contains(status, `"goroutines"`) {
        t.Errorf("Status response missing fields: %q", status)
    }
}
