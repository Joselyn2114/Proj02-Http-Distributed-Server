package core

import (
    "net"
    "testing"
)

// TestParseRequestSuccess cubre el caso exitoso de ParseRequest.
func TestParseRequestSuccess(t *testing.T) {
    raw := "POST /mypath?foo=bar HTTP/1.0\r\nX-Test: OK\r\nAnother: 123\r\n\r\n"
    req, err := ParseRequest(raw)
    if err != nil {
        t.Fatalf("Expected ParseRequest to succeed, got error: %v", err)
    }
    if req.Method != "POST" {
        t.Errorf("Expected method POST, got %s", req.Method)
    }
    if req.Target.Path != "/mypath" {
        t.Errorf("Expected path /mypath, got %s", req.Target.Path)
    }
    vals := req.Target.Query()
    if vals.Get("foo") != "bar" {
        t.Errorf("Expected query foo=bar, got %v", vals)
    }
    if req.Headers["X-Test"] != "OK" || req.Headers["Another"] != "123" {
        t.Errorf("Headers parsed incorrectly: %v", req.Headers)
    }
}

// TestReadRequestGetWithoutContentLength cubre la rama GET sin Content-Length.
func TestReadRequestGetWithoutContentLength(t *testing.T) {
    p1, p2 := net.Pipe()
    defer p2.Close()

    go func() {
        // Un GET sin Content-Length debe parsear y Body="" :contentReference[oaicite:0]{index=0}:contentReference[oaicite:1]{index=1}
        p1.Write([]byte("GET /hello HTTP/1.0\r\nFoo: Bar\r\n\r\n"))
        p1.Close()
    }()

    req, err := ReadRequest(p2)
    if err != nil {
        t.Fatalf("Expected no error for GET without body, got %v", err)
    }
    if req.Method != "GET" || req.Target.Path != "/hello" {
        t.Errorf("Unexpected request: %v %v", req.Method, req.Target)
    }
    if req.Body != "" {
        t.Errorf("Expected empty body, got %q", req.Body)
    }
}
