package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/todo-lists-app/todo-lists-api/internal/validate"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"

	"github.com/bugfixes/go-bugfixes/logs"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	healthcheck "github.com/keloran/go-healthcheck"
	probe "github.com/keloran/go-probe"
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
	go getNotifications(s.Config, errChan)

	return <-errChan
}

type injectData struct {
	Data string `json:"data"`
	IV   string `json:"iv"`
}

func getNotifications(cfg *config.Config, errChan chan error) {

}

//const cbSettings = gobreaker.Settings{
//	Name:        "HTTP Breaker",
//	Timeout:     time.Second * 60,
//	MaxRequests: 10,
//	Interval:    time.Minute * 1,
//	ReadyToTrip: func(counts gobreaker.Counts) bool {
//		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
//		return counts.Requests >= 10 && failureRatio >= 0.6
//	},
//}
//cb := gobreaker.NewCircuitBreaker(cbSettings)

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
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}).Handler)
	r.Get("/health", healthcheck.HTTP)
	r.Get("/probe", probe.HTTP)

	r.Route("/account", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")

			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			a := api.NewAccountService(r.Context(), *cfg, subject)
			account, err := a.GetAccount()
			if err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}

			if account == nil {
				account, err := a.CreateAccount()
				if err != nil {
					logs.Infof("Error: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					errChan <- err
					return
				}

				if err := AccountData(w, account); err != nil {
					logs.Infof("Error: %s", err)
					w.WriteHeader(http.StatusInternalServerError)
					errChan <- err
					return
				}
			}

			if err := AccountData(w, account); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				errChan <- err
				return
			}
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			logs.Local().Infof("Subject: %s", r.Header.Get("X-User-Subject"))
			w.Header().Set("debug", "delete account")
			logs.Local().Info("Delete account")
		})
	})

	r.Route("/notifications", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			type resp struct {
				ID       string `json:"id"`
				Title    string `json:"title"`
				Message  string `json:"message"`
				Read     bool   `json:"read"`
				Priority string `json:"priority"`
			}
			respJ := []resp{
				{
					ID:       "jyuuuooprf8",
					Title:    "test",
					Message:  "test message",
					Read:     false,
					Priority: "low",
				},
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respJ)
		})

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			n := api.NewNotificationService(r.Context(), *cfg, "test")
			if err := n.SendTestNotification(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logs.Infof("Error: %s", err)
			}
			w.WriteHeader(http.StatusOK)
		})
		r.Post("/subscribe", func(w http.ResponseWriter, r *http.Request) {
			subject := r.Header.Get("X-User-Subject")
			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v := validate.NewValidate(cfg, r.Context())
			valid, err := v.ValidateUser(subject)
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

			var sub webpush.Subscription
			if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
				logs.Infof("Error: %s", err)
			}

			n := api.NewNotificationService(r.Context(), *cfg, subject)
			if err := n.StoreUser(sub); err != nil {
				logs.Infof("Error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			logs.Infof("Subscription: %v", sub)
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

			v := validate.NewValidate(cfg, r.Context())
			valid, err := v.ValidateUser(subject)
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

			l := api.NewListService(r.Context(), *cfg, subject)
			list, err := l.GetList()
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					// logs.Infof("Error: %s", err)
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
			logs.Infof("Subject: %s", subject)

			if subject == "" {
				logs.Info("No Subject")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			v := validate.NewValidate(cfg, r.Context())
			valid, err := v.ValidateUser(subject)
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

			l := api.NewListService(r.Context(), *cfg, subject)
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

			v := validate.NewValidate(cfg, r.Context())
			valid, err := v.ValidateUser(subject)
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
			if _, err := api.NewListService(r.Context(), *cfg, subject).UpdateList(&api.StoredList{
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

			v := validate.NewValidate(cfg, r.Context())
			valid, err := v.ValidateUser(subject)
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
