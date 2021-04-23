package spec

import "github.com/zekroTJA/ranna/internal/models"

type Provider interface {
	Load() (models.SpecMap, error)
}
