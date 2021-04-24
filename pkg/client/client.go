package client

import "github.com/zekroTJA/ranna/pkg/models"

// Client provides an API endpoint wrapper
// for the ranna API.
type Client interface {

	// Spec requests the server's spec map.
	Spec() (spec models.SpecMap, err error)

	// Exec sends an executin request with the passed
	// execution parameters and returns either the
	// execution response or an error.
	//
	// An error response is only returned if the request
	// itself failed. If the executed code failed, this
	// will only be visible in the execution response.
	Exec(req models.ExecutionRequest) (res models.ExecutionResponse, err error)
}
