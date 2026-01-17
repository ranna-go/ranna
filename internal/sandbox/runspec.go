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
func (t RunSpec) GetAssembledHostDir() string {
	return path.Join(t.HostDir, t.Subdir)
}

// GetEntrypoint splits the entrypoint specification
// to a string array and returns it.
func (t RunSpec) GetEntrypoint() []string {
	return split(t.Entrypoint)
}

// GetCommandWithArgs splits the cmd specification
// (if passed) and appends given arguments.
func (t RunSpec) GetCommandWithArgs() []string {
	cmd := split(t.Cmd)
	return append(cmd, t.Arguments...)
}

// GetEnv assembles the environment variable map to
// a key-value string array.
//
// Also, the RANNA_HOSTDIR env variable is added here.
func (t RunSpec) GetEnv() (env []string) {
	if t.Environment == nil {
		t.Environment = make(map[string]string)
	}

	t.Environment["RANNA_HOSTDIR"] = t.HostDir

	env = make([]string, len(t.Environment))
	i := 0
	for k, v := range t.Environment {
		env[i] = fmt.Sprintf(`%s=%s`, k, v)
		i++
	}

	return
}

func split(v string) (res []string) {
	res = argRx.FindAllString(v, -1)
	for i, v := range res {
		res[i] = strings.ReplaceAll(v, "\"", "")
	}
	return
}
