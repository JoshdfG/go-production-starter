package main

import (
	"todo-clean/config"
	"todo-clean/internal/app"
	"todo-clean/pkg/logger"
)

// @title           Todo Clean API
// @version         1.0
// @BasePath        /v1
func main() {
	cfg, err := config.New()
	if err != nil {
		panic("config error: " + err.Error())
	}

	l := logger.New(cfg.Log.Level, cfg.App.Name, cfg.App.Version, cfg.App.Env)
	l.Info().Str("env", cfg.App.Env).Msgf("starting %s v%s", cfg.App.Name, cfg.App.Version)

	if err := app.Run(cfg, l); err != nil {
		l.Fatal().Err(err).Msg("app error")
	}
}
