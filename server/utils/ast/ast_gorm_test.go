package ast

import (
	"os"
	"path/filepath"
	"testing"
)

const A = 123

func TestAddRegisterTablesAst(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ast_test.go")
	content := []byte(`package ast

import (
	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/example"
)

func Register() {
	db := global.GetGlobalDBByDBName("test")
	db.AutoMigrate(example.ExaFile{})
}
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	AddRegisterTablesAst(path, "Register", "test", "testDB", "testModel")
}
