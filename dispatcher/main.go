package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

func mustAtoi(s string) int {
  i, _ := strconv.Atoi(s)
  return i
}

func main() {
  raw := os.Getenv("WORKERS") // "http://worker-pi:8081,http://worker-matrix:8082"
  workers := strings.Split(raw, ",")

  type state struct {
    URL   string `json:"url"`
    Alive bool   `json:"alive"`
    Load  int    `json:"load"`
    Done  int    `json:"done"`
  }
  states := make([]*state, len(workers))
  for i, url := range workers {
    states[i] = &state{URL: url}
  }

  srv := core.NewHttpServer()
  var idx uint64

  // 1. Health-check
  go func() {
    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
      for i, w := range workers {
        resp, err := http.Get(w + "/ping")
        if err != nil {
          states[i].Alive = false
          continue
        }
        var info struct {
          Status string `json:"status"`
          Load   int    `json:"load"`
          Done   int    `json:"done"`
        }
        _ = json.NewDecoder(resp.Body).Decode(&info)
        resp.Body.Close()
        states[i].Alive = true
        states[i].Load = info.Load
        states[i].Done = info.Done
      }
    }
  }()

  // 2. Proxy Round-Robin
  proxy := func(path string) core.Handle {
    return func(req *core.HttpRequest) (*core.HttpResponse, error) {
      n := uint64(len(workers))
      for a := uint64(0); a < n; a++ {
        i := int(atomic.AddUint64(&idx, 1) % n)
        if !states[i].Alive {
          continue
        }
				// build URL + raw query string from parsed URL
        target := workers[i] + path
        if q := req.Target.RawQuery; q != "" {
          target += "?" + q
        }
        newRequest, err := http.NewRequest(req.Method, target, strings.NewReader(req.Body))
        if err != nil {
          continue
        }
        newResponse, err := http.DefaultClient.Do(newRequest)
        if err != nil {
          continue
        }
        newResponseBody, err := io.ReadAll(newResponse.Body)
        defer newResponse.Body.Close()
        if err != nil {
          continue
        }
        newResponseHeaders := make(map[string]string)
        for k, v := range newResponse.Header {
          newResponseHeaders[k] = strings.Join(v, ", ")
        }
        return &core.HttpResponse{
          StatusCode: newResponse.StatusCode,
          StatusText: newResponse.Status,
          Headers:    newResponseHeaders,
          Body:       string(newResponseBody),
        }, nil
      }
			// if all workers dead
      return core.BadRequest().Text("todos los workers caídos"), nil
    }
  }

  for _, p := range []string{"/fibonacci", "/hash", "/simulate", "/pi", "/matrix"} {
    srv.Post(p, proxy(p))
  }

  // 3. /workers
  srv.Get("/workers", func(req *core.HttpRequest) (*core.HttpResponse, error) {
		// append nil error to match signature (*HttpResponse, error)
    return core.Ok().JsonObj(states), nil
  })

  // 4. Iniciar
  port := mustAtoi(os.Getenv("PORT"))
  srv.Start(port)
}
