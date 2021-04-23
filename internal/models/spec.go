package models

type SpecMap map[string]*Spec

type Spec struct {
	Image      string `json:"image" yaml:"image"`
	Entrypoint string `json:"entrypoint" yaml:"entrypoint"`
	Cmd        string `json:"cmd" yaml:"cmd"`
	Subdir     string `json:"subdir,omitempty" yaml:"subdir,omitempty"`
	HostDir    string `json:"hostdir,omitempty" yaml:"hostdir,omitempty"`
}
