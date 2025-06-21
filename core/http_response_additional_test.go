package core

import (
    "errors"
    "net"
    "testing"
    "time"
)

// ErrConn simula un net.Conn cuyo Write siempre falla.
type ErrConn struct{}

func (ErrConn) Read(_ []byte) (int, error)   { return 0, nil }
func (ErrConn) Write(_ []byte) (int, error)  { return 0, errors.New("boom") }
func (ErrConn) Close() error                 { return nil }
func (ErrConn) LocalAddr() net.Addr          { return dummyAddr{} }
func (ErrConn) RemoteAddr() net.Addr         { return dummyAddr{} }
func (ErrConn) SetDeadline(_ time.Time) error      { return nil }
func (ErrConn) SetReadDeadline(_ time.Time) error  { return nil }
func (ErrConn) SetWriteDeadline(_ time.Time) error { return nil }

type dummyAddr struct{}
func (dummyAddr) Network() string { return "tcp" }
func (dummyAddr) String() string  { return "addr" }

// TestJsonObjError cubre la rama de JsonObj cuando json.Marshal falla :contentReference[oaicite:2]{index=2}:contentReference[oaicite:3]{index=3}.
func TestJsonObjError(t *testing.T) {
    resp := Ok().JsonObj(make(chan int)) // un canal no es serializable
    if resp.StatusCode != 500 {
        t.Errorf("Expected status 500, got %d", resp.StatusCode)
    }
    if resp.StatusText != "Internal Server Error" {
        t.Errorf("Expected status text Internal Server Error, got %q", resp.StatusText)
    }
    if ct := resp.Headers["Content-Type"]; ct != "text/plain" {
        t.Errorf("Expected text/plain content type, got %q", ct)
    }
    if resp.Body != "json marshal error" {
        t.Errorf("Expected body \"json marshal error\", got %q", resp.Body)
    }
}

// TestWriteResponseError cubre la rama de error en WriteResponse :contentReference[oaicite:4]{index=4}:contentReference[oaicite:5]{index=5}.
func TestWriteResponseError(t *testing.T) {
    err := Ok().WriteResponse(ErrConn{})
    if err == nil || err.Error() != "boom" {
        t.Errorf("Expected WriteResponse to return boom, got %v", err)
    }
}
