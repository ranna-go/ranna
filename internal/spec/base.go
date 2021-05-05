package spec

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/ranna-go/ranna/pkg/models"
)

type baseProvider struct {
	m models.SpecMap
}

func newBaseProvider() *baseProvider {
	return &baseProvider{m: make(models.SpecMap)}
}

func (p *baseProvider) parseAndSet(data []byte, format string) (err error) {
	if strings.HasPrefix(format, ".") {
		format = format[1:]
	}

	var unmarshaller func([]byte, interface{}) error

	switch format {
	case "yml", "yaml", "application/x-yaml", "text/yaml":
		unmarshaller = yaml.Unmarshal
	case "json", "application/json", "text/json":
		unmarshaller = json.Unmarshal
	default:
		err = errors.New("unsupported file type")
		return
	}

	err = unmarshaller(data, &p.m)
	return
}

func (p *baseProvider) Spec() models.SpecMap {
	return p.m
}
