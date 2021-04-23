package spec

import "github.com/zekroTJA/ranna/internal/sandbox"

type SpecMap map[string]*sandbox.Spec

type Provider interface {
	Load() (SpecMap, error)
}
