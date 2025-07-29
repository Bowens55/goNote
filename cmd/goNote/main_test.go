package main

import (
	"testing"
)

// func loadDB(){

// }

func TestEnvFileExists(t *testing.T) {
	// Now you can test with controlled directories
	tempDir := t.TempDir()
	got, _ := EnvFileExists(tempDir)
	want := false

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
