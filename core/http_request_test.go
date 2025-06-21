package core

import (
	"fmt"
	"net"
	"testing"
)

func TestReadRequestPass(t *testing.T) {
	// Arrange
	conn1, conn2 := net.Pipe()

	defer conn2.Close()

	go func() {
		conn1.Write([]byte("GET / HTTP/1.0\r\nContent-Length: 7\r\nSkip\r\n: Skip\r\n\r\nContent"))
		conn1.Close()
	}()

	// Act
	request, err := ReadRequest(conn2)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, %v", err)
	}

	if request.Method != "GET" {
		t.Errorf("Expected method to be GET, not %s", request.Method)
	}

	if request.Target.Path != "/" {
		t.Errorf("Expected target path to be /, not %s", request.Target.Path)
	}

	if len(request.Headers) != 1 {
		t.Errorf("Expected 1 header, not %d", len(request.Headers))
	}

	if request.Headers["Content-Length"] != "7" {
		t.Errorf("Expected Content-Length to be 7, not %s", request.Headers["Content-Length"])
	}

	if request.Body != "Content" {
		t.Errorf("Expected body to be Content, not %s", request.Body)
	}
}

var RejectReadRequestTests = []string{
	// empty request
	"",
	// can't parse request
	"\r\n\r\n",
	// bad content length format
	"GET / HTTP/1.0\r\nContent-Length: A\r\n\r\nContent",
	// can't read body
	"GET / HTTP/1.0\r\nContent-Length: 7\r\n\r\n",
	// post request without content length
	"POST / HTTP/1.0\r\n\r\n",
	// bad method
	"BAD / HTTP/1.0\r\n\r\n",
	// bad version
	"GET / HTTP/0.0\r\n\r\n",
}

func TestReadRequestReject(t *testing.T) {
	for i, input := range RejectReadRequestTests {
		t.Run(fmt.Sprintf("TestReadRequestReject %d", i), func(t *testing.T) {
			// Arrange
			conn1, conn2 := net.Pipe()

			defer conn2.Close()

			go func() {
				conn1.Write([]byte(input))
				conn1.Close()
			}()

			// Act
			request, err := ReadRequest(conn2)

			// Assert
			if request != nil {
				t.Errorf("Expected no request, not %v", request)
			}

			if err == nil {
				t.Fatalf("Expected error")
			}
		})
	}
}

var RejectParseRequestTests = []string{
	// no header end
	"GET / HTTP/1.0\r\n",
	// no start line
	"\r\n\r\n",
	// no method or target
	"HTTP/1.0\r\n\r\n",
	// no target
	"GET\r\n\r\n",
	// bad target format
	"GET : HTTP/1.0\r\n\r\n",
}

func TestParseRequestReject(t *testing.T) {
	for i, input := range RejectParseRequestTests {
		t.Run(fmt.Sprintf("TestParseRequestReject %d", i), func(t *testing.T) {
			// Act
			request, err := ParseRequest(input)

			// Assert
			if request != nil {
				t.Errorf("Expected no request, not %v", request)
			}

			if err == nil {
				t.Fatalf("Expected error")
			}
		})
	}
}
