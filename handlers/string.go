package handlers

import (
    "crypto/sha256"
    "encoding/hex"
    "github.com/KateGF/Http-Server-Project-SO/core"
    "strings"
)

// helper: extrae query param
func getParam(req *core.HttpRequest, key string) (string, *core.HttpResponse) {
    val := req.Target.Query().Get(key)
    if val == "" {
        return "", core.BadRequest().Text(key + " is required")
    }
    return val, nil
}

// /reverse?text=…
func ReverseHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    text, errResp := getParam(req, "text")
    if errResp != nil {
        return errResp, nil
    }
    runes := []rune(text)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return core.Ok().Text(string(runes)), nil
}

// /toupper?text=…
func ToUpperHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    text, errResp := getParam(req, "text")
    if errResp != nil {
        return errResp, nil
    }
    return core.Ok().Text(strings.ToUpper(text)), nil
}

// /hash?text=…
func HashHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    text, errResp := getParam(req, "text")
    if errResp != nil {
        return errResp, nil
    }
    sum := sha256.Sum256([]byte(text))
    return core.Ok().Text(hex.EncodeToString(sum[:])), nil
}

// / (ruta raíz)
func RootHandler(req *core.HttpRequest) (*core.HttpResponse, error) {
    return core.Ok().Text(`Servidor HTTP activo. Rutas disponibles:
GET  /reverse?text=...
GET  /toupper?text=...
GET  /hash?text=...
GET  /timestamp
GET  /random?count=n&min=a&max=b
GET  /simulate?seconds=s&task=name
GET  /sleep?seconds=s
GET  /loadtest?tasks=n&sleep=s
GET  /status
GET  /help`), nil
}
