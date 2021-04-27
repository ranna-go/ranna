package spec

import "github.com/ranna-go/ranna/pkg/models"

type Provider interface {
	Load() error
	Spec() models.SpecMap
}
