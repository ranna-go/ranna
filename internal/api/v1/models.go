package v1

type executionRequest struct {
	Language    string            `json:"language"`
	Code        string            `json:"code"`
	Arguments   []string          `json:"arguments"`
	Environment map[string]string `json:"environment"`
}

type executionResponse struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}
