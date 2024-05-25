package sandbox

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/ranna-go/ranna/pkg/models"
)

var argRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)

// RunSpec wraps a spec and extends runtime
// information like arguments and environment variables
// passed to the sandbox as well as the sub directory and
// host dir used to inject the code snippet into
// the sandbox.
type RunSpec struct {
	models.Spec

	Arguments   []string          `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Subdir      string            `json:"subdir,omitempty" yaml:"subdir,omitempty"`
	HostDir     string            `json:"hostdir,omitempty" yaml:"hostdir,omitempty"`
}

// GetAssembledHostDir returns the joined directory
// of host dir and sub dir.
func (s RunSpec) GetAssembledHostDir() string {
	return path.Join(s.HostDir, s.Subdir)
}

// GetEntrypoint splits the entrypoint specification
// to a string array and returns it.
func (s RunSpec) GetEntrypoint() []string {
	return split(s.Entrypoint)
}

// GetCommandWithArgs splits the cmd specification
// (if passed) and appends given arguments.
func (s RunSpec) GetCommandWithArgs() []string {
	cmd := split(s.Cmd)
	return append(cmd, s.Arguments...)
}

// GetEnv assembles the environment variable map to
// a key-value string array.
//
// Also, the RANNA_HOSTDIR env variable is added here.
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
