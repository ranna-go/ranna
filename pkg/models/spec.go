package models

import "strings"

// Spec defines a code environment specification.
type Spec struct {
	Image      string `json:"image,omitempty" yaml:"image,omitempty"`
	Entrypoint string `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	FileName   string `json:"filename,omitempty" yaml:"filename,omitempty"`
	Cmd        string `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Registry   string `json:"registry,omitempty" yaml:"registry,omitempty"`
	Use        string `json:"use,omitempty" yaml:"use,omitempty"`
	Example    string `json:"example,omitempty" yaml:"example,omitempty"`
	Template   string `json:"template,omitempty" yaml:"template,omitempty"`
}

// SupportsTemplating checks if a spec supports templating for inline expressions
func (spec *Spec) SupportsTemplating() bool {
	return spec.Template != "" && strings.Contains(spec.Template, "$${CODE}")
}

// SpecMap wraps a map[string]*Spec.
type SpecMap map[string]*Spec
