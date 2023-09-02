package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/keloran/go-healthcheck"
	"github.com/keloran/go-probe"
	validate "github.com/todo-lists-app/go-validate-user"
	"github.com/todo-lists-app/todo-lists-api/internal/api"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
)

// Service is the service
type Service struct {
	Config *config.Config
}

// Start the service
func (s *Service) Start() error {
	errChan := make(chan error)

	go startHTTP(s.Config, errChan)

	return <-errChan
}

type injectData struct {
	Data string `json:"data"`
	IV   string `json:"iv"`
}

//golint:ignore(gocyclo)
func startHTTP(cfg *config.Config, errChan chan error) {
	p := fmt.Sprintf(":%d", cfg.Local.HTTPPort)
	logs.Local().Infof("starting http on %s", p)

	allowedOrigins := []string{
		"http://localhost:3000",
		"https://todo-list.app",
		"https://beta.todo-list.app",
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
			"X-User-Subject",
			"X-User-Access-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)
	r.Get("/health", healthcheck.HTTP)
	r.Get("/probe", probe.HTTP)

	r.Route("/account", func(r chi.Router) {
		//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		//	subject := r.Header.Get("X-User-Subject")
		//
		//	if subject == "" {
		//		logs.Info("No Subject")
		//		w.WriteHeader(http.StatusUnauthorized)
		//		return
		//	}
		//
		//	a := api.NewAccountService(r.Context(), *cfg, subject)
		//	account, err := a.GetAccount()
		//	if err != nil {
		//		logs.Infof("Error: %s", err)
		//		w.WriteHeader(http.StatusInternalServerError)
		//		errChan <- err
		//		return
		//	}
		//
		//	if account == nil {
		//		account, err := a.CreateAccount()
		//		if err != nil {
		//			logs.Infof("Error: %s", err)
		//			w.WriteHeader(http.StatusInternalServerError)
		//			errChan <- err
		//			return
		//		}
		//
		//		if err := AccountData(w, account); err != nil {
		//			logs.Infof("Error: %s", err)
		//			w.WriteHeader(http.StatusInternalServerError)
		//			errChan <- err
		//			return
		//		}
		//	}
		//
		//	if err := AccountData(w, account); err != nil {
		//		logs.Infof("Error: %s", err)
		//		w.WriteHeader(http.StatusInternalServerError)
		//		errChan <- err
		//		return
		//	}
		//})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("X-User-Access-Token")
			if accessToken == "" {
				logs.Info("No Access Token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v, err := validate.NewValidate(r.Context(), cfg.Services.Identity, cfg.Local.Development).GetClient()
			if err != nil {
				logs.Infof("validate client err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			valid, err := v.ValidateUser(accessToken, subject)
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !valid {
				logs.Info("invalid user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			a, err := api.NewAccountService(r.Context(), *cfg, subject, accessToken).GetClient()
			if err != nil {
				logs.Infof("Error Get Account Client: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := a.DeleteAccount(); err != nil {
				logs.Infof("Error Delete Account: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		})
	})

	r.Route("/list", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			accessToken := r.Header.Get("X-User-Access-Token")
			if accessToken == "" {
				logs.Info("No Access Token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v, err := validate.NewValidate(r.Context(), cfg.Services.Identity, cfg.Local.Development).GetClient()
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			valid, err := v.ValidateUser(accessToken, subject)
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !valid {
				logs.Infof("invalid user: %+v", r.Header)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			l, err := api.NewListService(r.Context(), *cfg, subject).GetClient()
			if err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			list, err := l.GetList()
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					w.WriteHeader(http.StatusOK)
					return
				}

				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			if list == nil {
				if err := NoLists(w); err != nil {
					logs.Infof("Error: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					errChan <- err
					return
				}
				return
			}

			if err := ListExists(w, list); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			accessToken := r.Header.Get("X-User-Access-Token")
			if accessToken == "" {
				logs.Info("No Access Token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v, err := validate.NewValidate(r.Context(), cfg.Services.Identity, cfg.Local.Development).GetClient()
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			valid, err := v.ValidateUser(accessToken, subject)
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !valid {
				logs.Info("invalid user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			id := injectData{}
			if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			l, err := api.NewListService(r.Context(), *cfg, subject).GetClient()
			if err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			stored, err := l.CreateList(&api.StoredList{
				UserID: subject,
				Data:   id.Data,
				IV:     id.IV,
			})
			if err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			if err := ListExists(w, stored); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			accessToken := r.Header.Get("X-User-Access-Token")
			if accessToken == "" {
				logs.Info("No Access Token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v, err := validate.NewValidate(r.Context(), cfg.Services.Identity, cfg.Local.Development).GetClient()
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			valid, err := v.ValidateUser(accessToken, subject)
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !valid {
				logs.Info("invalid user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			id := injectData{}
			if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			l, err := api.NewListService(r.Context(), *cfg, subject).GetClient()
			if err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			if _, err := l.UpdateList(&api.StoredList{
				UserID: subject,
				Data:   id.Data,
				IV:     id.IV,
			}); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			w.Header().Set("debug", "put list")
			w.WriteHeader(http.StatusOK)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			accessToken := r.Header.Get("X-User-Access-Token")
			if accessToken == "" {
				logs.Info("No Access Token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v, err := validate.NewValidate(r.Context(), cfg.Services.Identity, cfg.Local.Development).GetClient()
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			valid, err := v.ValidateUser(accessToken, subject)
			if err != nil {
				logs.Infof("validate user err: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !valid {
				logs.Info("invalid user")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			logs.Infof("Subject: %s", subject)
			w.Header().Set("debug", "delete list")
			w.WriteHeader(http.StatusNotImplemented)
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
