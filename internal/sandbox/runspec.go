package sandbox

import (
	"fmt"
	"path"
	"strings"

	"github.com/zekroTJA/ranna/pkg/models"
)

type RunSpec struct {
	models.Spec

	Arguments   []string          `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Subdir      string            `json:"subdir,omitempty" yaml:"subdir,omitempty"`
	HostDir     string            `json:"hostdir,omitempty" yaml:"hostdir,omitempty"`
}

func (s RunSpec) GetAssambledHostDir() string {
	return path.Join(s.HostDir, s.Subdir)
}

func (s RunSpec) GetEntrypoint() []string {
	return strings.Split(s.Entrypoint, " ")
}

func (s RunSpec) GetCommandWithArgs() []string {
	cmd := strings.Split(s.Cmd, " ")
	return append(cmd, s.Arguments...)
}

func (s RunSpec) GetEnv() (env []string) {
	if s.Environment != nil {
		env = make([]string, len(s.Environment))
		i := 0
		for k, v := range s.Environment {
			env[i] = fmt.Sprintf(`%s=%s`, k, v)
			i++
		}
	}

	return
}
