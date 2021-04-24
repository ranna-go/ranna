// Package models provides all models used by
// the ranna API.
package models

// ErrorModel is the reponse model returned from
// the ranna API when something went wrong.
type ErrorModel struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Context string `json:"context,omitempty"`
}
