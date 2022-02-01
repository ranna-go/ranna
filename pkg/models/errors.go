package models

import (
	"fmt"
)

var (
	ErrInvalidMessageType = WsError{400, "invalid message type"}
	ErrInvalidOpCode      = WsError{400, "invalid operation code"}
	ErrEmptyCode          = WsError{400, "code is empty"}
	ErrUnimplemented      = WsError{501, "this action is currently unimplemented"}
	ErrSandboxNotRunning  = WsError{400, "sandbox is not running"}
	ErrRateLimited        = WsError{429, "you have been rate limited"}
)

type WsError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e WsError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
