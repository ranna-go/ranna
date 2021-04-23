package spec

import "github.com/zekroTJA/ranna/internal/models"

type Provider interface {
	Load() error
	Spec() models.SpecMap
}
