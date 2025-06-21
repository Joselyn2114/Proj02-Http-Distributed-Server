package handlers

import (
	"github.com/KateGF/Http-Server-Project-SO/core"
	"net/url"
	"testing"
)

func makeReq(method, raw string) *core.HttpRequest {
	u, _ := url.Parse(raw)
	return core.NewHttpRequest(method, u, map[string]string{}, "")
}

func TestReverse(t *testing.T) {
	res, _ := ReverseHandler(makeReq("GET", "/reverse?text=hola"))
	if res.Body != "aloh" {
		t.Errorf("Reverse incorrecto: %q", res.Body)
	}
	res2, _ := ReverseHandler(makeReq("GET", "/reverse"))
	if res2.StatusCode != 400 {
		t.Errorf("esperaba BadRequest si falta text")
	}
}

func TestToUpper(t *testing.T) {
	res, _ := ToUpperHandler(makeReq("GET", "/toupper?text=GoLang"))
	if res.Body != "GOLANG" {
		t.Errorf("ToUpper incorrecto: %q", res.Body)
	}
}

func TestHash(t *testing.T) {
	res, _ := HashHandler(makeReq("GET", "/hash?text=abc"))
	// SHA256("abc") = "..."
	want := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"

	if res.Body != want {
		t.Errorf("Hash incorrecto: got %q, want %q", res.Body, want)
	}
}

func TestRoot(t *testing.T) {
	res, _ := RootHandler(makeReq("GET", "/"))
	expected := "Servidor HTTP activo. Rutas disponibles:\n" +
		"GET  /reverse?text=...\n" +
		"GET  /toupper?text=...\n" +
		"GET  /hash?text=...\n" +
		"GET  /timestamp\n" +
		"GET  /random?count=n&min=a&max=b\n" +
		"GET  /simulate?seconds=s&task=name\n" +
		"GET  /sleep?seconds=s\n" +
		"GET  /loadtest?tasks=n&sleep=s\n" +
		"GET  /status\n" +
		"GET  /help"
	if res.Body != expected {
		t.Errorf("RootHandler incorrecto: got %q, want %q", res.Body, expected)
	}
}
