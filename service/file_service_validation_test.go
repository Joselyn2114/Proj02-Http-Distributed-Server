// service/file_service_validation_test.go
package service

import "testing"

func TestCreateFileValidation(t *testing.T) {
    // Nombre vacío → error
    err := CreateFile("", "content", 1)
    if err == nil {
        t.Errorf("Expected error for empty filename")
    }

    // Repeat negativo → error
    err = CreateFile("name.txt", "c", -1)
    if err == nil {
        t.Errorf("Expected error for negative repeat")
    }
}

func TestDeleteFileValidation(t *testing.T) {
    // Nombre vacío → error
    err := DeleteFile("")
    if err == nil {
        t.Errorf("Expected error for empty filename")
    }
}
