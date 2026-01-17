package models

type EventCode int

const (
	EventPong EventCode = iota
	EventError
	EventSpawn
	EventLog
	EventStop
)

type OpCode int

const (
	OpPing OpCode = iota
	OpExec
	OpKill
)

type Event struct {
	Code  EventCode `json:"code"`
	Nonce int       `json:"nonce,omitempty"`
	Data  any       `json:"data,omitempty"`
}

type Operation struct {
	Op    OpCode `json:"op"`
	Nonce int    `json:"nonce"`
}

type OperationExec struct {
	Operation
	Args ExecutionRequest `json:"args"`
}

type OperationKill struct {
	Operation
	Args DataRunId `json:"args"`
}

type DataRunId struct {
	RunId string `json:"runid"`
}

type DataLog struct {
	DataRunId
	StdOut string `json:"stdout,omitempty"`
	StdErr string `json:"stderr,omitempty"`
}

type DataStop struct {
	DataRunId
	ExecTimeMS int `json:"exectimems"`
}

type DataError struct {
	DataRunId
	Error error
}

type DataSpawn struct {
	DataRunId
}
