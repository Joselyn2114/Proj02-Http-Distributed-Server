package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/KateGF/Http-Server-Project-SO/worker/matrix"
)

// workerMatrixHandler is the HTTP handler for our test workers. It performs the partial multiplication.
func workerMatrixHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		A [][]float64 `json:"a"`
		B [][]float64 `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad JSON data: "+err.Error(), http.StatusBadRequest)
		return
	}

	matrixA, err := matrix.NewMatrix(data.A)
	if err != nil {
		http.Error(w, "Invalid matrix 'a': "+err.Error(), http.StatusBadRequest)
		return
	}

	matrixB, err := matrix.NewMatrix(data.B)
	if err != nil {
		http.Error(w, "Invalid matrix 'b': "+err.Error(), http.StatusBadRequest)
		return
	}

	resultMatrix, err := matrixA.Multiply(matrixB)
	if err != nil {
		http.Error(w, "Matrix multiplication error: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	jsonData, err := resultMatrix.ToJson()
	if err != nil {
		http.Error(w, "Error serializing result matrix: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(jsonData))
}

// createAndStartWorker sets up and runs a single worker HTTP server on a random available port.
func createAndStartWorker(t *testing.T) *http.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/matrix/part", workerMatrixHandler)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen on a port for worker: %v", err)
	}

	server := &http.Server{
		Addr:    listener.Addr().String(),
		Handler: mux,
	}

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Logf("Worker server error: %v", err)
		}
	}()

	return server
}

func TestMatrixMultiplicationIntegration(t *testing.T) {
	// Lightweight type to decode worker status without mutex
	type workerStatus struct {
		URL       string `json:"url"`
		Active    bool   `json:"active"`
		TasksDone int    `json:"tasks_done"`
	}

	// Arrange

	// Reset dispatcher's global state
	mu.Lock()
	workers = nil
	rrIndex = 0
	mu.Unlock()

	// Start the Dispatcher Server
	dispatcherMux := http.NewServeMux()
	dispatcherMux.HandleFunc("/register", RegisterHandler)
	dispatcherMux.HandleFunc("/matrix", MatrixHandler)
	dispatcherMux.HandleFunc("/workers", StatusHandler)

	dispatcherListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener for dispatcher: %v", err)
	}
	dispatcherServer := &http.Server{Handler: dispatcherMux}
	go func() {
		if err := dispatcherServer.Serve(dispatcherListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Logf("Dispatcher server error: %v", err)
		}
	}()
	defer dispatcherServer.Shutdown(context.Background())
	dispatcherURL := "http://" + dispatcherListener.Addr().String()

	// Start and register Worker Servers
	numWorkers := 3
	workerServers := make([]*http.Server, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workerServers[i] = createAndStartWorker(t)
		defer workerServers[i].Shutdown(context.Background())

		workerURL := "http://" + workerServers[i].Addr
		data, _ := json.Marshal(map[string]string{"URL": workerURL})
		resp, err := http.Post(dispatcherURL+"/register", "application/json", bytes.NewReader(data))
		if err != nil {
			t.Fatalf("Failed to register worker %s: %v", workerURL, err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Expected status 204 for worker registration, got %d", resp.StatusCode)
		}
		resp.Body.Close()
	}

	// Wait for all workers to become active
	var activeCount int
	for i := 0; i < 20; i++ { // Timeout after ~2 seconds
		resp, err := http.Get(dispatcherURL + "/workers")
		if err != nil {
			t.Fatalf("Failed to get worker status: %v", err)
		}
		var statusList []workerStatus
		if err := json.NewDecoder(resp.Body).Decode(&statusList); err != nil {
			t.Fatalf("Failed to decode worker status: %v", err)
		}
		resp.Body.Close()

		activeCount = 0
		for _, s := range statusList {
			if s.Active {
				activeCount++
			}
		}
		if activeCount == numWorkers {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if activeCount != numWorkers {
		t.Fatalf("Timed out waiting for all workers to become active. Found %d active workers", activeCount)
	}
	t.Logf("All %d workers successfully registered and are active.", numWorkers)

	// Define matrix inputs, expected output, and the HTTP request
	matrixA := [][]float64{
		{7, 91, 56, 17, 36, 68, 77, 85, 96, 45},
		{32, 94, 9, 3, 20, 11, 48, 74, 18, 59},
		{65, 87, 81, 71, 79, 99, 13, 8, 38, 50},
		{93, 29, 64, 21, 5, 27, 40, 19, 70, 46},
		{10, 62, 53, 49, 97, 88, 33, 76, 2, 41},
		{52, 60, 31, 89, 4, 26, 90, 14, 28, 78},
		{23, 72, 1, 98, 66, 12, 42, 57, 80, 35},
		{84, 15, 6, 25, 61, 63, 75, 47, 95, 30},
		{100, 24, 86, 43, 69, 58, 22, 92, 34, 55},
		{16, 82, 44, 73, 51, 67, 54, 39, 90, 9},
	}
	matrixB := [][]float64{
		{80, 4, 30, 71, 46, 66, 73, 54, 99, 36},
		{24, 91, 10, 93, 85, 34, 19, 48, 68, 79},
		{55, 60, 67, 3, 29, 95, 12, 87, 10, 22},
		{27, 50, 45, 96, 62, 44, 88, 74, 53, 10},
		{38, 2, 6, 8, 1, 92, 28, 49, 13, 59},
		{64, 40, 5, 52, 35, 78, 11, 84, 72, 97},
		{90, 81, 14, 21, 69, 15, 31, 63, 43, 6},
		{83, 77, 89, 100, 16, 57, 20, 65, 8, 42},
		{98, 5, 41, 70, 76, 60, 82, 18, 94, 25},
		{9, 17, 37, 23, 40, 8, 4, 9, 7, 7},
	}
	expectedResult := [][]float64{
		{35801, 29338, 20437, 32456, 28920, 30360, 18303, 30861, 27036, 24310},
		{19613, 20531, 12992, 23086, 18538, 15483, 9855, 17838, 15622, 15312},
		{29006, 23414, 16713, 30428, 26593, 36388, 20798, 34647, 29624, 27683},
		{26592, 14826, 15301, 21650, 21075, 22945, 17937, 21478, 23715, 13281},
		{25687, 24258, 16523, 26067, 18767, 31022, 13411, 31322, 18432, 24900},
		{24232, 22860, 14936, 26337, 26505, 19181, 19216, 25137, 23339, 13316},
		{26211, 21002, 16579, 31879, 25059, 24340, 22805, 24385, 24126, 17467},
		{34666, 16632, 15116, 27156, 23476, 27994, 22265, 25669, 28791, 19310},
		{34244, 22323, 23566, 30593, 21673, 36404, 20672, 33942, 25007, 22746},
		{30863, 24578, 16424, 31335, 27454, 29659, 21713, 29795, 28113, 22535},
	}
	requestData, _ := json.Marshal(map[string][][]float64{
		"a": matrixA,
		"b": matrixB,
	})
	req, err := http.NewRequest("POST", dispatcherURL+"/matrix", bytes.NewReader(requestData))
	if err != nil {
		t.Fatalf("Failed to create matrix request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Act

	// Send the multiplication request to the dispatcher
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send matrix request to dispatcher: %v", err)
	}
	defer resp.Body.Close()

	// Assert

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200 OK from dispatcher, got %d. Body: %s", resp.StatusCode, string(body))
	}

	// Decode the response body
	var actualResult [][]float64
	if err := json.NewDecoder(resp.Body).Decode(&actualResult); err != nil {
		t.Fatalf("Failed to decode the result from dispatcher: %v", err)
	}

	// Compare the actual result with the expected result
	if !reflect.DeepEqual(expectedResult, actualResult) {
		t.Errorf("Matrix multiplication result is incorrect.\nExpected:\n%v\nGot:\n%v", expectedResult, actualResult)
	} else {
		t.Logf("Matrix multiplication successful. Result matches expected output.")
	}
}
