package models

// Spec defines a code environment specification.
type Spec struct {
	Image      string `json:"image" yaml:"image"`
	Entrypoint string `json:"entrypoint" yaml:"entrypoint"`
	FileName   string `json:"filename" yaml:"filename"`
}

// SpecMap wraps a map[string]*Spec.
type SpecMap map[string]*Spec

func (m SpecMap) Get(key string) (s Spec, ok bool) {
	sp, ok := m[key]
	if !ok {
		return
	}

	s = *sp
	return
}
