package ast

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KeSilent/study-hub/server/global"
)

func TestRollRouterBack(t *testing.T) {
	root := t.TempDir()
	initDir := filepath.Join(root, "initialize")
	if err := os.MkdirAll(initDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	routerPath := filepath.Join(initDir, "router.go")
	routerContent := []byte(`package initialize

func Routers() {
	{
		tttRouter := 1
		_ = tttRouter
		tttRouter.InitTesttttRouter(PrivateGroup)
	}
}
`)
	if err := os.WriteFile(routerPath, routerContent, 0o644); err != nil {
		t.Fatalf("write router.go: %v", err)
	}

	global.GVA_CONFIG.AutoCode.Root = root
	global.GVA_CONFIG.AutoCode.Server = ""

	RollRouterBack("ttt", "Testttt")
}

func TestRollGormBack(t *testing.T) {
	root := t.TempDir()
	initDir := filepath.Join(root, "initialize")
	if err := os.MkdirAll(initDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	gormPath := filepath.Join(initDir, "gorm.go")
	gormContent := []byte(`package initialize

import "github.com/KeSilent/study-hub/server/model/ttt"

func RegisterTables() {
	db.AutoMigrate(ttt.Testttt{})
}
`)
	if err := os.WriteFile(gormPath, gormContent, 0o644); err != nil {
		t.Fatalf("write gorm.go: %v", err)
	}

	global.GVA_CONFIG.AutoCode.Root = root
	global.GVA_CONFIG.AutoCode.Server = ""

	RollGormBack("ttt", "Testttt")
}
