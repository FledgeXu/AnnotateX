package auth

import (
	"embed"

	"fmt"
	"github.com/casbin/casbin/v2"
	casbin_fs_adapter "github.com/naucon/casbin-fs-adapter"
)

//go:embed model.conf policy.csv
var embedFiles embed.FS
var Enforcer *casbin.Enforcer

func InitCasbinEnforcer() {
	model, err := casbin_fs_adapter.NewModel(embedFiles, "model.conf")
	if err != nil {
		panic(fmt.Errorf("failed to load Casbin model: %w", err))
	}
	policies := casbin_fs_adapter.NewAdapter(embedFiles, "policy.csv")
	Enforcer, err = casbin.NewEnforcer(model, policies)
	if err != nil {
		panic(fmt.Errorf("failed to load Casbin policy: %w", err))
	}
}
