package main

import (
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"github.com/todo-lists-app/todo-lists-api/internal/service"
)

var (
	BuildVersion = "dev"
	BuildHash    = "unknown"
	ServiceName  = "base-service"
)

func main() {
	logs.Local().Info(fmt.Sprintf("Starting %s", ServiceName))
	logs.Local().Info(fmt.Sprintf("Version: %s, Hash: %s", BuildVersion, BuildHash))

	cfg, err := config.Build()
	if err != nil {
		_ = logs.Local().Errorf("config: %v", err)
		return
	}

	s := &service.Service{
		Config: cfg,
	}

	if err := s.Start(); err != nil {
		_ = logs.Local().Errorf("start service: %v", err)
		return
	}
}
