package ast

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddRouterCode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "router_test.go")
	content := []byte(`package ast

func Routers() {
	{
		// placeholder
	}
	{
		// placeholder 2
	}
}
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	AddRouterCode(path, "Routers", "testRouter", "GVAStruct")
}
