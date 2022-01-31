package ws

import "errors"

var (
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrInvalidOpCode      = errors.New("invalid operation code")
	ErrUnimplemented      = errors.New("this action is currently unimplemented")
	ErrEmptyCode          = errors.New("code is empty")
	ErrSandboxNotRunning  = errors.New("sandbox is not running")
	ErrRateLimited        = errors.New("you have been rate limited")
)
