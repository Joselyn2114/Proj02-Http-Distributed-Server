// integration_test.go
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

// rawReq abre un socket, envía un request HTTP/1.0 válido (con Content-Length: 0 si aplica),
// y devuelve únicamente el body de la respuesta.
func rawReq(t *testing.T, method, path string) string {
	t.Helper()
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	// Monta la petición
	if method == "POST" || method == "DELETE" {
		fmt.Fprintf(conn, "%s %s HTTP/1.0\r\nContent-Length: 0\r\n\r\n", method, path)
	} else {
		fmt.Fprintf(conn, "%s %s HTTP/1.0\r\n\r\n", method, path)
	}

	// Lee toda la respuesta
	respData, err := io.ReadAll(conn)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	// Separa headers y body
	parts := strings.SplitN(string(respData), "\r\n\r\n", 2)
	if len(parts) < 2 {
		t.Fatalf("respuesta inválida:\n%s", respData)
	}
	return parts[1]
}

func startServer(t *testing.T) {
	t.Helper()
	go main() // arranca tu main.go
	time.Sleep(100 * time.Millisecond)
}

func TestIntegrationEndpoints(t *testing.T) {
	startServer(t)

	// 1) Endpoints de consulta
	checks := []struct {
		method, path, want string
	}{
		{"GET", "/fibonacci?num=7", "13"},
		{"GET", "/reverse?text=hola", "aloh"},
		{"GET", "/toupper?text=GoLang", "GOLANG"},
		{"GET", "/hash?text=abc", "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
	}
	for _, c := range checks {
		body := rawReq(t, c.method, c.path)
		if body != c.want {
			t.Errorf("%s %s: got %q, want %q", c.method, c.path, body, c.want)
		}
	}

	// 2) Endpoints de archivos
	fname := "integration_test.txt"
	defer os.Remove(fname)

	// createfile
	createPath := fmt.Sprintf("/createfile?name=%s&content=X&repeat=4", fname)
	body := rawReq(t, "POST", createPath)
	if !strings.Contains(body, "File created successfully") {
		t.Errorf("CreateFile: unexpected body %q", body)
	}
	data, _ := os.ReadFile(fname)
	if string(data) != "XXXX" {
		t.Errorf("CreateFile content: got %q, want %q", data, "XXXX")
	}

	// deletefile
	deletePath := fmt.Sprintf("/deletefile?name=%s", fname)
	body = rawReq(t, "DELETE", deletePath)
	if !strings.Contains(body, "File deleted successfully") {
		t.Errorf("DeleteFile: unexpected body %q", body)
	}
	if _, err := os.Stat(fname); !os.IsNotExist(err) {
		t.Errorf("DeleteFile: file still exists")
	}
}
