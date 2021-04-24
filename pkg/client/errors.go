package client

import (
	"fmt"
	"net/http"

	"github.com/zekroTJA/ranna/pkg/models"
)

type ResponseError struct {
	*models.ErrorModel
	Response *http.Response
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Error)
}
