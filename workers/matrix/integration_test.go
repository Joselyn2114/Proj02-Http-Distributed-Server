package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/KateGF/Http-Server-Project-SO/core"
)

func initMatrixServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/matrix", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("(Matrix) ← got %s %s from %s",
			r.Method, r.URL.RequestURI(), r.RemoteAddr,
		)

		bodyBytes, _ := io.ReadAll(r.Body)
		coreReq := &core.HttpRequest{Body: string(bodyBytes)}
		resp, _ := MatrixHandler(coreReq)

		log.Printf("(Matrix) → responding %d to %s",
			resp.StatusCode, r.RemoteAddr,
		)

		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(resp.Body))
	})

	return mux
}

func initDispatcher(matrixURL string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/matrix", func(w http.ResponseWriter, r *http.Request) {
		target := matrixURL + "/matrix"
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}

		log.Printf("(Dispatcher) → proxying %s %s from %s to %s",
			r.Method, r.URL.RequestURI(), r.RemoteAddr, matrixURL,
		)

		bodyBytes, _ := io.ReadAll(r.Body)
		newReq, _ := http.NewRequest(r.Method, target, bytes.NewBuffer(bodyBytes))
		newReq.Header = r.Header.Clone()

		newResp, err := http.DefaultClient.Do(newReq)
		if err != nil {
			http.Error(w, "proxy error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer newResp.Body.Close()

		log.Printf("(Dispatcher) ← got %d from %s, forwarding to %s",
			newResp.StatusCode, matrixURL, r.RemoteAddr,
		)

		w.WriteHeader(newResp.StatusCode)
		respBody, _ := io.ReadAll(newResp.Body)
		w.Write(respBody)
	})

	return mux
}

func TestIntegration_DispatcherMatrix(t *testing.T) {
	matrixSrv := httptest.NewServer(initMatrixServer())
	defer matrixSrv.Close()
	t.Logf("matrix server started at %s", matrixSrv.URL)

	dispSrv := httptest.NewServer(initDispatcher(matrixSrv.URL))
	defer dispSrv.Close()
	t.Logf("dispatcher server started at %s", dispSrv.URL)

	payload := map[string][][]float64{
		"A": {{1, 2}, {3, 4}},
		"B": {{5, 6}, {7, 8}},
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(dispSrv.URL+"/matrix", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var result [][]float64
	json.Unmarshal(body, &result)

	expected := [][]float64{{19, 22}, {43, 50}}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
