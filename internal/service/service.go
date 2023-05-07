package service

import (
	"encoding/json"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	healthcheck "github.com/keloran/go-healthcheck"
	probe "github.com/keloran/go-probe"
	"github.com/todo-lists-app/todo-lists-api/internal/api"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"net/http"
	"time"
)

type Service struct {
	Config *config.Config
}

func (s *Service) Start() error {
	errChan := make(chan error)

	go startHTTP(s.Config, errChan)

	return <-errChan
}

func NoLists(w http.ResponseWriter) error {
	type NoList struct {
		Message string       `json:"message"`
		Data    api.TodoList `json:"data"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(NoList{
		Message: "No Lists",
		Data:    api.TodoList{},
	})
}

func ListExists(w http.ResponseWriter, l *api.TodoList) error {
	type List struct {
		Message string       `json:"message"`
		Data    api.TodoList `json:"data"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(List{
		Data: *l,
	})
}

func startHTTP(cfg *config.Config, errChan chan error) {
	p := fmt.Sprintf(":%d", cfg.Local.HTTPPort)
	logs.Local().Infof("starting http on %s", p)

	allowedOrigins := []string{
		"http://localhost:3000",
		"https://todo-list.app",
	}
	if cfg.Local.Development {
		allowedOrigins = append(allowedOrigins, "http://*")
	}

	r := chi.NewRouter()
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-User-Token",
			"X-User-Subject",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)
	r.Get("/health", healthcheck.HTTP)
	r.Get("/probe", probe.HTTP)

	r.Route("/list", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			token := r.Header.Get("X-User-Token")

			if subject == "" || token == "" {
				logs.Local().Info("No Subject or Token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			l := api.NewListService(r.Context(), *cfg, subject, token)
			list, err := l.GetList()
			if err != nil {
				logs.Local().Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if list == nil {
				if err := NoLists(w); err != nil {
					logs.Local().Infof("Error: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				return
			}

			if err := ListExists(w, list); err != nil {
				logs.Local().Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			logs.Local().Infof("Subject: %s", r.Header.Get("X-User-Subject"))
			w.Header().Set("debug", "post list")
			logs.Local().Info("Post List")
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			logs.Local().Infof("Subject: %s", r.Header.Get("X-User-Subject"))
			w.Header().Set("debug", "put list")
			logs.Local().Info("Put List")
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			logs.Local().Infof("Subject: %s", r.Header.Get("X-User-Subject"))
			w.Header().Set("debug", "delete list")
			logs.Local().Info("Delete List")
		})
	})

	srv := &http.Server{
		Addr:              p,
		Handler:           r,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		errChan <- err
		return
	}
}
