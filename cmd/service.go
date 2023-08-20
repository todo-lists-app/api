// Package main run the app
package main

import (
	"fmt"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"github.com/todo-lists-app/todo-lists-api/internal/service"
)

var (
	// BuildVersion is the version of the app
	BuildVersion = "dev"
	// BuildHash is the git-hash of the app
	BuildHash = "unknown"
	// ServiceName is the name of the service
	ServiceName = "base-service"
)

func runApp() error {
	logs.Local().Info(fmt.Sprintf("Starting %s", ServiceName))
	logs.Local().Info(fmt.Sprintf("Version: %s, Hash: %s", BuildVersion, BuildHash))

	cfg, err := config.Build()
	if err != nil {
		return logs.Local().Errorf("config: %v", err)
	}

	s := &service.Service{
		Config: cfg,
	}

	if err := s.Start(); err != nil {
		return logs.Local().Errorf("start service: %v", err)
	}

	return nil
}

func main() {
	if err := runApp(); err != nil {
		_ = logs.Local().Errorf("run app: %v", err)
	}
}
