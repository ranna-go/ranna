package sandbox

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/ranna-go/ranna/pkg/models"
)

var argRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)

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
	return split(s.Entrypoint)
}

func (s RunSpec) GetCommandWithArgs() []string {
	cmd := split(s.Cmd)
	return append(cmd, s.Arguments...)
}

func (s RunSpec) GetEnv() (env []string) {
	if s.Environment == nil {
		s.Environment = make(map[string]string)
	}

	s.Environment["RANNA_HOSTDIR"] = s.HostDir

	env = make([]string, len(s.Environment))
	i := 0
	for k, v := range s.Environment {
		env[i] = fmt.Sprintf(`%s=%s`, k, v)
		i++
	}

	return
}

func split(v string) (res []string) {
	res = argRx.FindAllString(v, -1)
	for i, v := range res {
		res[i] = strings.Replace(v, "\"", "", -1)
	}
	return
}
