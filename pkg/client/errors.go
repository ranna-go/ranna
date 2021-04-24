package client

import (
	"fmt"
	"net/http"

	"github.com/zekroTJA/ranna/pkg/models"
)

// ResponseError is an error which wraps
// a response ErrorModel and the Response
// object reference itself.
type ResponseError struct {
	ErrorModel *models.ErrorModel
	Response   *http.Response
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.ErrorModel.Code, e.ErrorModel.Error)
}
