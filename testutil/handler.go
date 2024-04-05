package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("cannot unmarshal got %q: %v", got, err)
	}
	if diff := cmp.Diff(jg, jw); diff != "" {
		t.Errorf("got differs: (-got + want)\n%s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	assert := assert.New(t)

	t.Helper()
	t.Cleanup(func() { got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(got.StatusCode, status)

	if len(gb) == 0 && len(body) == 0 {
		return
	}

	AssertJSON(t, body, gb)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}
