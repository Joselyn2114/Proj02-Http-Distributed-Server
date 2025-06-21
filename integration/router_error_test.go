// integration/router_error_test.go
package integration

import (
    "bufio"
    "fmt"
    "net"
    "strings"
    "testing"
)

// sendStatusLine envía la petición y devuelve únicamente la línea de estado HTTP.
func sendStatusLine(t *testing.T, req string) string {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        t.Fatalf("Dial failed: %v", err)
    }
    defer conn.Close()

    fmt.Fprint(conn, req)
    reader := bufio.NewReader(conn)
    line, err := reader.ReadString('\n')
    if err != nil {
        t.Fatalf("Read status line failed: %v", err)
    }
    return line
}

func TestNotFoundRoute(t *testing.T) {
    status := sendStatusLine(t,
        "GET /no_such_route HTTP/1.0\r\nHost: test\r\n\r\n",
    )
    if !strings.Contains(status, "404") {
        t.Errorf("Expected 404 Not Found, got %q", status)
    }
}

func TestBadMethod(t *testing.T) {
    status := sendStatusLine(t,
        "POST /fibonacci?num=5 HTTP/1.0\r\nHost: test\r\n\r\n",
    )
    if !strings.Contains(status, "400") {
        t.Errorf("Expected 400 Bad Request, got %q", status)
    }
}
