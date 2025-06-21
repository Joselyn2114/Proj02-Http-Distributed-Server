package core

import (
	"net"
	"testing"
)

func TestHttpResponse(t *testing.T) {
	// Act
	response := Ok()

	response.SetStatusCode(418)
	response.SetStatusText("I'm a teapot")
	response.SetBody("Content")

	message := response.String()

	// Assert
	if response.StatusCode != 418 {
		t.Errorf("Expected status code to be 418, not %d", response.StatusCode)
	}

	if response.StatusText != "I'm a teapot" {
		t.Errorf("Expected status text to be I'm a teapot, not %s", response.StatusText)
	}

	if response.Body != "Content" {
		t.Errorf("Expected body to be Content, not %s", response.Body)
	}

	expected := "HTTP/1.0 418 I'm a teapot\r\nContent-Length: 7\r\n\r\nContent"

	if message != expected {
		t.Errorf("Expected message to be %s, not %s", expected, message)
	}
}

func TestHttpResponseNotFound(t *testing.T) {
	// Act
	response := NotFound()

	// Assert
	if response.StatusCode != 404 {
		t.Errorf("Expected status code to be 404, not %d", response.StatusCode)
	}

	if response.StatusText != "Not Found" {
		t.Errorf("Expected status text to be Not Found, not %s", response.StatusText)
	}

	if response.Body != "" {
		t.Errorf("Expected body to be empty, not %s", response.Body)
	}

	expected := "HTTP/1.0 404 Not Found\r\nContent-Length: 0\r\n\r\n"

	if response.String() != expected {
		t.Errorf("Expected message to be %s, not %s", expected, response.String())
	}
}

func TestHttpResponseText(t *testing.T) {
	// Act
	response := Ok().Text("Text")

	// Assert
	if response.Body != "Text" {
		t.Errorf("Expected body to be Text, not %s", response.Body)
	}

	expected := "HTTP/1.0 200 OK\r\nContent-Length: 4\r\nContent-Type: text/plain\r\n\r\nText"

	if response.String() != expected {
		t.Errorf("Expected message to be %s, not %s", expected, response.String())
	}
}

func TestHttpResponseJson(t *testing.T) {
	// Act
	response := Ok().Json("{\"k\": \"v\"}")

	// Assert
	if response.Body != "{\"k\": \"v\"}" {
		t.Errorf("Expected body to be {\"k\": \"v\"}, not %s", response.Body)
	}

	expected := "HTTP/1.0 200 OK\r\nContent-Length: 10\r\nContent-Type: application/json\r\n\r\n{\"k\": \"v\"}"

	if response.String() != expected {
		t.Errorf("Expected message to be %s, not %s", expected, response.String())
	}
}

func TestHttpResponseWriteResponse(t *testing.T) {
	// Arrange
	conn1, conn2 := net.Pipe()

	defer conn2.Close()

	// Act
	go func() {
		response := Ok()
		response.WriteResponse(conn1)
		conn1.Close()
	}()

	// Assert
	buffer := make([]byte, 1024)
	n, err := conn2.Read(buffer)

	if err != nil {
		t.Errorf("Expected no error, %v", err)
	}

	message := string(buffer[:n])

	expected := "HTTP/1.0 200 OK\r\nContent-Length: 0\r\n\r\n"

	if message != expected {
		t.Errorf("Expected message to be %s, not %s", expected, message)
	}
}
