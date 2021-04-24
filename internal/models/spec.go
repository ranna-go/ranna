package models

import (
	"path"
	"strings"
)

type SpecMap map[string]*Spec

func (m SpecMap) Get(key string) (s Spec, ok bool) {
	sp, ok := m[key]
	if !ok {
		return
	}

	s = *sp
	return
}

type Spec struct {
	Image      string `json:"image" yaml:"image"`
	Entrypoint string `json:"entrypoint" yaml:"entrypoint"`
	FileName   string `json:"filename" yaml:"filename"`

	Cmd         string            `json:"cmd" yaml:"cmd"`
	Arguments   []string          `json:"arguments" yaml:"arguments"`
	Environment map[string]string `json:"environment" yaml:"environment"`
	Subdir      string            `json:"subdir,omitempty" yaml:"subdir,omitempty"`
	HostDir     string            `json:"hostdir,omitempty" yaml:"hostdir,omitempty"`
}

func (s Spec) GetAssambledHostDir() string {
	return path.Join(s.HostDir, s.Subdir)
}

func (s Spec) GetEntrypoint() []string {
	return strings.Split(s.Entrypoint, " ")
}

func (s Spec) GetCommandWithArgs() []string {
	cmd := strings.Split(s.Cmd, " ")
	return append(cmd, s.Arguments...)
}
