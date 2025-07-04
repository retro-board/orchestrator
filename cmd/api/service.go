package main

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	ConfigBuilder "github.com/keloran/go-config"
	"github.com/retro-board/orchestrator/internal"
)

var (
	BuildVersion = "0.0.1"
	BuildHash    = "unknown"
	ServiceName  = "service"
)

type ProjectConfig struct{}

func (pc ProjectConfig) Build(cfg *ConfigBuilder.Config) error {
	type PC struct {
		// Railway
		RailwayPort string `env:"PORT" envDefault:"3000"`
		OnRailway   bool   `env:"ON_RAILWAY" envDefault:"false"`
	}
	p := PC{}

	if err := env.Parse(&p); err != nil {
		return logs.Errorf("Failed to parse services: %v", err)
	}
	if cfg.ProjectProperties == nil {
		cfg.ProjectProperties = make(map[string]interface{})
	}
	cfg.ProjectProperties["railway_port"] = p.RailwayPort
	cfg.ProjectProperties["on_railway"] = p.OnRailway

	return nil
}

func main() {
	logs.Logf("Starting %s version %s (build %s)", ServiceName, BuildVersion, BuildHash)
	c := ConfigBuilder.NewConfigNoVault()

	if err := c.Build(
		ConfigBuilder.Local,
		ConfigBuilder.Bugfixes,
		ConfigBuilder.Postgres,
		ConfigBuilder.Flags,
		ConfigBuilder.WithProjectConfigurator(ProjectConfig{})); err != nil {
		logs.Fatalf("Failed to build config: %v", err)
	}

	if err := internal.New(c).Start(); err != nil {
		logs.Fatalf("Failed to start service: %v", err)
	}
}
