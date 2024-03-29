package models

import (
	"regexp"
	"strings"
)

// Spec defines a code environment specification.
type Spec struct {
	Image      string      `json:"image,omitempty" yaml:"image,omitempty"`
	Entrypoint string      `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
	FileName   string      `json:"filename,omitempty" yaml:"filename,omitempty"`
	Cmd        string      `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Registry   string      `json:"registry,omitempty" yaml:"registry,omitempty"`
	Use        string      `json:"use,omitempty" yaml:"use,omitempty"`
	Language   string      `json:"language,omitempty" yaml:"language,omitempty"`
	Example    string      `json:"example,omitempty" yaml:"example,omitempty"`
	Inline     *InlineSpec `json:"inline,omitempty" yaml:"inline,omitempty"`
}

type InlineSpec struct {
	ImportRegex         string         `json:"import_regex" yaml:"import_regex"`
	ImportRegexCompiled *regexp.Regexp `json:"-" yaml:"-"`
	Template            string         `json:"template" yaml:"template"`
}

// SupportsTemplating checks if a spec supports templating for inline expressions
func (spec *Spec) SupportsTemplating() bool {
	return spec.Inline != nil &&
		strings.Contains(spec.Inline.Template, "$${CODE}")
}

// SpecMap wraps a map[string]*Spec.
type SpecMap map[string]*Spec

// Get returns the last spec in the map
// also following "use" pointers.
func (m SpecMap) Get(key string) (s *Spec, ok bool) {
	for {
		s = m[key]
		if s == nil {
			return
		}
		if s.Use == "" {
			ok = true
			return
		}
		key = s.Use
	}
}
