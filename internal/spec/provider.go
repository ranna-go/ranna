package spec

import "github.com/zekroTJA/ranna/pkg/models"

type Provider interface {
	Load() error
	Spec() models.SpecMap
}
