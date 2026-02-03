package ast

import (
	"os"
	"path/filepath"
	"testing"
)

func TestImportForAutoEnter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "enter.go")
	content := []byte(`package test

type ApiGroup struct {
}
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	ImportForAutoEnter(path, "ApiGroup", "test")
}
