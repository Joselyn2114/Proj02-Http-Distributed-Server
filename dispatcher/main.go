package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type WorkerInfo struct {
    URL       string
    Active    bool
    TasksDone int
    mu        sync.Mutex
}

var (
    workers []*WorkerInfo
    rrIndex int
    mu      sync.Mutex
)

// --- Registro dinámico de workers ---

// RegisterHandler añade un nuevo worker al dispatcher, sin duplicados.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct{ URL string }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.URL == "" {
        http.Error(w, "Bad JSON", http.StatusBadRequest)
        return
    }
    mu.Lock()
    defer mu.Unlock()
    // Si ya existe, simplemente lo reactivamos si estaba inactivo
    for _, wk := range workers {
        if wk.URL == payload.URL {
            wk.mu.Lock()
            wk.Active = true
            wk.mu.Unlock()
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    // Si no existe, lo añadimos al slice
    workers = append(workers, &WorkerInfo{URL: payload.URL, Active: true})
    w.WriteHeader(http.StatusNoContent)
}

// UnregisterHandler elimina un worker cuando apaga
func UnregisterHandler(w http.ResponseWriter, r *http.Request) {
    var payload struct{ URL string }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.URL == "" {
        http.Error(w, "Bad JSON", http.StatusBadRequest)
        return
    }
    mu.Lock()
    defer mu.Unlock()
    for i, wk := range workers {
        if wk.URL == payload.URL {
            workers = append(workers[:i], workers[i+1:]...)
            break
        }
    }
    w.WriteHeader(http.StatusNoContent)
}

// --- Health-checker periódico ---

func HealthChecker() {
    for {
        mu.Lock()
        pool := make([]*WorkerInfo, len(workers))
        copy(pool, workers)
        mu.Unlock()

        var wg sync.WaitGroup
        for _, wk := range pool {
            wg.Add(1)
            go func(wk *WorkerInfo) {
                defer wg.Done()
                client := http.Client{Timeout: 2 * time.Second}
                resp, err := client.Get(wk.URL + "/ping")
                wk.mu.Lock()
                wk.Active = (err == nil && resp.StatusCode == http.StatusOK)
                wk.mu.Unlock()
            }(wk)
        }
        wg.Wait()
        time.Sleep(5 * time.Second)
    }
}

// --- Utilidades de workers ---

// GetNextWorker devuelve el siguiente activo en round-robin
func GetNextWorker() *WorkerInfo {
    mu.Lock()
    defer mu.Unlock()
    n := len(workers)
    for i := 0; i < n; i++ {
        rrIndex = (rrIndex + 1) % n
        if workers[rrIndex].Active {
            return workers[rrIndex]
        }
    }
    return nil
}

// GetActiveWorkers devuelve la lista de workers actualmente activos
func GetActiveWorkers() []*WorkerInfo {
    mu.Lock()
    defer mu.Unlock()
    active := make([]*WorkerInfo, 0, len(workers))
    for _, wk := range workers {
        wk.mu.Lock()
        if wk.Active {
            active = append(active, wk)
        }
        wk.mu.Unlock()
    }
    return active
}

// DoRequestWithRetry intenta hasta maxTries repartir la petición si un worker falla
func DoRequestWithRetry(method, url string, payload []byte, headers http.Header, maxTries int) (*http.Response, error) {
    var lastErr error
    tried := make(map[string]bool)

    for attempt := 0; attempt < maxTries; attempt++ {
        wk := GetNextWorker()
        if wk == nil {
            break
        }
        // evita reintentar el mismo worker inmediatamente
        if tried[wk.URL] {
            continue
        }
        tried[wk.URL] = true

        req, err := http.NewRequest(method, wk.URL+url, bytes.NewReader(payload))
        if err != nil {
            lastErr = err
            continue
        }
        req.Header = headers.Clone()
        req.Header.Del("Transfer-Encoding")
        req.Header.Set("Content-Length", strconv.Itoa(len(payload)))

        resp, err := (&http.Client{}).Do(req)
        if err == nil && resp.StatusCode < 500 {
            wk.mu.Lock()
            wk.TasksDone++
            wk.mu.Unlock()
            return resp, nil
        }
        // marcar inactivo y guardar error
        wk.mu.Lock()
        wk.Active = false
        wk.mu.Unlock()
        if err != nil {
            lastErr = err
        } else {
            lastErr = errors.New("status " + resp.Status)
            resp.Body.Close()
        }
    }
    return nil, fmt.Errorf("all workers failed: %v", lastErr)
}

// ProxyHandler reenvía cualquier ruta GENÉRICA a un worker con retry
func ProxyHandler(w http.ResponseWriter, r *http.Request) {
    payload, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error leyendo cuerpo", http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    resp, err := DoRequestWithRetry(r.Method, r.RequestURI, payload, r.Header, len(workers))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    for k, vs := range resp.Header {
        for _, v := range vs {
            w.Header().Add(k, v)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

// StatusHandler muestra el estado actual de todos los workers
func StatusHandler(w http.ResponseWriter, _ *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    out := make([]map[string]interface{}, 0, len(workers))
    for _, wk := range workers {
        wk.mu.Lock()
        out = append(out, map[string]interface{}{
            "url":        wk.URL,
            "active":     wk.Active,
            "tasks_done": wk.TasksDone,
        })
        wk.mu.Unlock()
    }
    json.NewEncoder(w).Encode(out)
}

// --- Endpoint /matrix: split, distribuir, merge ---

// SplitMatrixRows divide A en n bloques de filas
func SplitMatrixRows(A [][]float64, n int) [][][]float64 {
    m := len(A)
    size := (m + n - 1) / n
    parts := make([][][]float64, 0, n)
    for i := 0; i < m; i += size {
        end := i + size
        if end > m {
            end = m
        }
        parts = append(parts, A[i:end])
    }
    return parts
}

// StitchMatrix recompone bloques en una sola matriz
func StitchMatrix(parts [][][]float64) [][]float64 {
    var result [][]float64
    for _, block := range parts {
        result = append(result, block...)
    }
    return result
}

func MatrixHandler(w http.ResponseWriter, r *http.Request) {
    // 1) Decode del JSON de entrada
    var payload struct {
        A [][]float64 `json:"a"`
        B [][]float64 `json:"b"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // 2) Split de A en bloques de filas
    rowBlocks := SplitMatrixRows(payload.A, len(GetActiveWorkers()))

    // 3) Preparar slice donde guardaremos cada respuesta
    responses := make([][][]float64, len(rowBlocks))

    // 4) Lanzar un goroutine por bloque
    var wg sync.WaitGroup
    wg.Add(len(rowBlocks))
    for i, blk := range rowBlocks {
        go func(idx int, blockChunk [][]float64) {
            defer wg.Done()

            // serializar sub-payload
            subPayload, _ := json.Marshal(map[string]any{
                "a": blockChunk,
                "b": payload.B,
            })

            // hacer POST con retry
            resp, err := DoRequestWithRetry(
                "POST",
                "/matrix/part",
                subPayload,
                http.Header{"Content-Type": []string{"application/json"}},
                len(workers),
            )
            if err != nil {
                log.Printf("block %d failed: %v", idx, err)
                return
            }
            defer resp.Body.Close()

            // decodificar respuesta del worker
            var partRes [][]float64
            if err := json.NewDecoder(resp.Body).Decode(&partRes); err != nil {
                log.Printf("decode block %d failed: %v", idx, err)
                return
            }
            responses[idx] = partRes
        }(i, blk)
    }
    wg.Wait()

    // 5) Stitch de las sub-matrices y respuesta final
    result := StitchMatrix(responses)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func main() {
    // Arranque estático inicial
    // initWorkers := []string{
    //     "http://worker1:8080",
    //     "http://worker2:8080",
    //     "http://worker3:8080",
    // }
    // for _, u := range initWorkers {
    //     workers = append(workers, &WorkerInfo{URL: u, Active: true})
    // }

    go HealthChecker()

    http.HandleFunc("/register", RegisterHandler)
    http.HandleFunc("/unregister", UnregisterHandler)
    http.HandleFunc("/workers", StatusHandler)
    http.HandleFunc("/matrix", MatrixHandler)    // endpoint completo
    http.HandleFunc("/", ProxyHandler)           // proxy para todo lo demás

    log.Println("Dispatcher escuchando en :8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
