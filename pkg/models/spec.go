package models

// Spec defines a code environment specification.
type Spec struct {
	Image      string `json:"image" yaml:"image"`
	Entrypoint string `json:"entrypoint" yaml:"entrypoint"`
	FileName   string `json:"filename" yaml:"filename"`
	Cmd        string `json:"cmd,omitempty" yaml:"cmd,omitempty"`
	Registry   string `json:"registry,omitempty" yaml:"registry,omitempty"`
	Use        string `json:"use,omitempty" yaml:"use,omitempty"`
}

// SpecMap wraps a map[string]*Spec.
type SpecMap map[string]*Spec

func (m SpecMap) get(key string, isAlias bool) (s Spec, ok bool) {
	sp, ok := m[key]
	if !ok {
		return
	}

	if sp.Use != "" {
		if isAlias {
			ok = false
			return
		}
		return m.get(sp.Use, true)
	}

	s = *sp
	return
}

func (m SpecMap) Get(key string) (Spec, bool) {
	return m.get(key, false)
}
