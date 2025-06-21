package service

import (
	"fmt"
	"github.com/KateGF/Http-Server-Project-SO/core"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func Clean(t *testing.T, filename string) {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, filename)
	os.Remove(path)
}

func TestCreateFileHandler(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
		wantErr  bool
		before   func()
	}{
		{"temp/test.txt", "name=temp/test.txt&content=A&repeat=10", "AAAAAAAAAA", false, nil},
		{"temp/test.txt", "content=A&repeat=10", "", true, nil},
		{"temp/test.txt", "name=temp/test.txt&repeat=10", "", true, nil},
		{"temp/test.txt", "name=temp/test.txt&content=A", "", true, nil},
		{"temp/test.txt", "name=temp/test.txt&content=A&repeat=A", "", true, nil},
		{"temp/test.txt", "name=temp/test.txt&content=A&repeat=-10", "", true, nil},
		{"temp/test.txt", "name=temp/test.txt&content=A&repeat=10", "AAAAAAAAAA", true, func() {
			os.WriteFile("temp/test.txt", []byte("test"), os.ModePerm)
		}},
		{"../temp/test.txt", "name=../test.txt&content=A&repeat=10", "", true, nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestCreateFileHandler %d", i), func(t *testing.T) {
			// Before
			Clean(t, tt.name)
			if tt.before != nil {
				tt.before()
			}
			defer Clean(t, tt.name)

			// Arrange
			target, _ := url.Parse(fmt.Sprintf("/createfile?%s", tt.query))
			request := core.NewHttpRequest("POST", target, map[string]string{}, "")

			// Act
			response, _ := CreateFileHandler(request)

			// Assert
			if !tt.wantErr && response.StatusCode != 200 {
				t.Fatalf("Expected status code to be 200, not %d", response.StatusCode)
			}

			if tt.wantErr && response.StatusCode == 200 {
				t.Fatalf("Expected status code to be not 200")
			}

			if !tt.wantErr {
				if _, err := os.Stat(tt.name); os.IsNotExist(err) {
					t.Fatalf("Expected file %s to exist ", tt.name)
				}

				content, err := os.ReadFile(tt.name)
				if err != nil {
					t.Fatalf("Can't read file %s: %v", tt.name, err)
				}

				contentStr := string(content)
				if contentStr != tt.expected {
					t.Errorf("Expected content to be %s, not %s", tt.expected, contentStr)
				}
			}
		})
	}
}

func TestDeleteFileHandler(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
		before  func()
	}{
		{"temp/test.txt", "name=temp/test.txt", false, func() {
			os.WriteFile("temp/test.txt", []byte("test"), os.ModePerm)
		}},
		{"temp/test.txt", "", true, nil},
		{"temp/test.txt", "name=temp/test.txt", true, nil},
		{"../temp/test.txt", "name=../test.txt", true, nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestDeleteFileHandler %d", i), func(t *testing.T) {
			// Before
			Clean(t, tt.name)
			if tt.before != nil {
				tt.before()
			}
			defer Clean(t, tt.name)

			// Arrange
			target, _ := url.Parse(fmt.Sprintf("/deletefile?%s", tt.query))
			request := core.NewHttpRequest("DELETE", target, map[string]string{}, "")

			// Act
			response, _ := DeleteFileHandler(request)

			// Assert
			if !tt.wantErr && response.StatusCode != 200 {
				t.Fatalf("Expected status code to be 200, not %d", response.StatusCode)
			}

			if tt.wantErr && response.StatusCode == 200 {
				t.Fatalf("Expected status code to be not 200")
			}

			if !tt.wantErr {
				if _, err := os.Stat(tt.name); !os.IsNotExist(err) {
					t.Fatalf("Expected file %s to not exist ", tt.name)
				}
			}
		})
	}
}
