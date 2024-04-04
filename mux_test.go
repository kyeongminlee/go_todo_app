package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMux(t *testing.T) {
	assert := assert.New(t)

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut := NewMux()
	sut.ServeHTTP(writer, request)
	response := writer.Result()
	t.Cleanup(func() { response.Body.Close() })

	assert.Equal(response.StatusCode, http.StatusOK)

	got, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "ok"}`
	assert.Equal(string(got), want)
}
