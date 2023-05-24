// Package api is the service itself
package api

import (
	"context"

	"github.com/todo-lists-app/todo-lists-api/internal/config"
)

// AccountDetails is the account details
type AccountDetails struct {
	Salt string `json:"salt,omitempty"`
}

// Account is the account service
type Account struct {
	context.Context
	config.Config
	Subject string
}

// NewAccountService creates a new account service
func NewAccountService(ctx context.Context, cfg config.Config, subject string) *Account {
	return &Account{
		Context: ctx,
		Config:  cfg,
		Subject: subject,
	}
}

// GetAccount gets an account for the user
func (a *Account) GetAccount() (*AccountDetails, error) {
	return &AccountDetails{
		Salt: "salt",
	}, nil
}

// CreateAccount creates an account for the user
func (a *Account) CreateAccount() (*AccountDetails, error) {
	return &AccountDetails{
		Salt: "salt",
	}, nil
}
