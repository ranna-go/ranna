package client

import "github.com/zekroTJA/ranna/pkg/models"

type Client interface {
	Spec() (spec models.SpecMap, err error)
	Exec(req models.ExecutionRequest) (res models.ExecutionResponse, err error)
}
