package handlers

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// sample struct to decode into
type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestGetBody_Success(t *testing.T) {
	body := `{"name": "Alice", "age": 30}`
	r := &http.Request{
		Body: io.NopCloser(strings.NewReader(body)),
	}

	var result testStruct
	err := getBody(r, &result)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if result.Name != "Alice" || result.Age != 30 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGetBody_InvalidJSON(t *testing.T) {
	body := `{"name": "Bob", "age": }` // invalid JSON
	r := &http.Request{
		Body: io.NopCloser(strings.NewReader(body)),
	}

	var result testStruct
	err := getBody(r, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetBody_ReadError(t *testing.T) {
	r := &http.Request{
		Body: io.NopCloser(&errReader{}),
	}

	var result testStruct
	err := getBody(r, &result)
	if err == nil {
		t.Fatal("expected read error, got nil")
	}
}

// Custom reader that always fails
type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
