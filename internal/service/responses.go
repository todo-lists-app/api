// Package service provides the service layer for the todo-lists-api.
package service

import (
	"encoding/json"
	"net/http"

	"github.com/todo-lists-app/todo-lists-api/internal/api"
)

// NoLists returns a 200 with no lists.
func NoLists(w http.ResponseWriter) error {
	type NoList struct {
		Message string         `json:"message,omitempty"`
		Data    api.StoredList `json:"data,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(NoList{
		Message: "No Lists",
		Data:    api.StoredList{},
	})
}

// ListExists returns the list data for the user.
func ListExists(w http.ResponseWriter, l *api.StoredList) error {
	type List struct {
		Message string `json:"message,omitempty"`
		Data    string `json:"data,omitempty"`
		IV      string `json:"iv,omitempty"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(List{
		Data: l.Data,
		IV:   l.IV,
	})
}

// AccountData returns the account details for the user.
//func AccountData(w http.ResponseWriter, a *api.AccountDetails) error {
//	type Account struct {
//		Message string             `json:"message,omitempty"`
//		Data    api.AccountDetails `json:"data,omitempty"`
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	return json.NewEncoder(w).Encode(Account{
//		Data: *a,
//	})
//}
