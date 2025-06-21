package core

import (
	"fmt"
	"math/rand/v2"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestGetPost(t *testing.T) {
	// Arrange
	server := NewHttpServer()

	handler := func(request *HttpRequest) (*HttpResponse, error) {
		return Ok(), nil
	}

	// Act
	server.Get("/get", handler)

	server.Post("/post", handler)

	// Assert
	if len(server.Handlers) != 2 {
		t.Fatalf("Expected 2 handlers, not %d", len(server.Handlers))
	}

	if server.Handlers[0].Method != "GET" || server.Handlers[0].Path != "/get" {
		t.Errorf("Expected GET handler to be added")
	}

	if server.Handlers[1].Method != "POST" || server.Handlers[1].Path != "/post" {
		t.Errorf("Expected POST handler to be added")
	}
}

func TestSortHandlers(t *testing.T) {
	// Arrange
	server := NewHttpServer()

	handler := func(request *HttpRequest) (*HttpResponse, error) {
		return Ok(), nil
	}

	server.AddHandler("GET", "/", handler)
	server.AddHandler("GET", "/users", handler)
	server.AddHandler("GET", "/users/name/posts", handler)
	server.AddHandler("GET", "/users/name", handler)

	expected := []string{"/users/name/posts", "/users/name", "/users", "/"}

	// Act
	server.SortHandlers()

	// Assert
	order := make([]string, len(server.Handlers))
	for i, h := range server.Handlers {
		order[i] = h.Path
	}

	if !reflect.DeepEqual(order, expected) {
		t.Errorf("Expected handlers to be sorted as %v, not %v", expected, order)
	}
}

var MatchPathTests = []struct {
	requestPath string
	handlerPath string
	expected    bool
}{
	{"/", "/", true},
	{"/users", "/", false},
	{"/users", "/users", true},
	{"/users/name", "/users", true},
	{"/users/name/posts", "/users", true},
	{"/users", "/users/name", false},
	{"/users", "/users/name/posts", false},
	{"/", "/users", false},
}

func TestMatchPath(t *testing.T) {
	for i, test := range MatchPathTests {
		t.Run(fmt.Sprintf("TestMatchPath %d", i), func(t *testing.T) {
			// Act
			act := MatchPath(test.requestPath, test.handlerPath)

			// Assert
			if act != test.expected {
				t.Errorf("Expected %v, not %v", test.expected, act)
			}
		})
	}
}

func TestStart(t *testing.T) {
	server := NewHttpServer()

	server.Get("/get", func(request *HttpRequest) (*HttpResponse, error) {
		return Ok(), nil
	})

	server.Get("/post", func(request *HttpRequest) (*HttpResponse, error) {
		return Ok(), nil
	})

	server.Delete("/delete", func(request *HttpRequest) (*HttpResponse, error) {
		return Ok(), nil
	})

	port := rand.IntN(1000) + 8080

	go func() {
		err := server.Start(port)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}()

	ready := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)

		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		defer conn.Close()

		fmt.Fprintf(conn, "GET /get HTTP/1.0\r\nHost: localhost:%d\r\n\r\n", port)

		close(ready)
	}()

	<-ready

	server.Stop()
}

func TestStartNotFound(t *testing.T) {
	server := NewHttpServer()

	port := rand.IntN(1000) + 8080

	go func() {
		err := server.Start(port)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}()

	ready := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)

		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		defer conn.Close()

		fmt.Fprintf(conn, "GET /get HTTP/1.0\r\nHost: localhost:%d\r\n\r\n", port)

		close(ready)
	}()

	<-ready

	server.Stop()
}

func TestStartUnknownError(t *testing.T) {
	server := NewHttpServer()

	port := rand.IntN(1000) + 8080

	server.Get("/get", func(request *HttpRequest) (*HttpResponse, error) {
		return nil, fmt.Errorf("Unknown error")
	})

	go func() {
		err := server.Start(port)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
	}()

	ready := make(chan int)

	go func() {
		time.Sleep(1 * time.Second)

		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		defer conn.Close()

		fmt.Fprintf(conn, "GET /get HTTP/1.0\r\nHost: localhost:%d\r\n\r\n", port)

		close(ready)
	}()

	<-ready

	server.Stop()
}
