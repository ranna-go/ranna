package models

// ExecutionRequest is the execution
// request model.
type ExecutionRequest struct {
	Language    string            `json:"language"`
	Code        string            `json:"code"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

// ExecutionResponse is the response
// model received on execution request.
type ExecutionResponse struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}
