package models

type ExecutionRequest struct {
	Language    string            `json:"language"`
	Code        string            `json:"code"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

type ExecutionResponse struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}
