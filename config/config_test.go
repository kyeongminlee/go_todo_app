package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("Port", fmt.Sprint(wantPort))

	got, err := New()
	assert.NoError(t, err, "cannot create config")
	// if err != nil {
	// 	t.Fatalf("cannot create config: %v", err)
	// }

	assert.Equal(t, wantPort, got.Port, "Port value mismatch")
	// if got.Port != wantPort {
	// 	t.Errorf("want %d, but %d", wantPort, got.Port)
	// }

	wantEnv := "dev"
	assert.Equal(t, wantEnv, got.Env, "Env value mismatch")
	// if got.Env != wantEnv {
	// 	t.Errorf("want %s, but %s", wantEnv, got.Env)
	// }
}
