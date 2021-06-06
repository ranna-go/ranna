package spec

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/ranna-go/ranna/pkg/models"
)

// baseProvider provides a common base for all providers
// to store the spec map as SafeSpecMap, expose it and
// parse text data to a spec map.
type baseProvider struct {
	m *SafeSpecMap
}

// newBaseProvider initializes a new baseProvider
func newBaseProvider() *baseProvider {
	return &baseProvider{m: nil}
}

// parseAndSet takes a spec definition as text data and
// a format (either format name or MIME type) to be used
// to parse the given data.
//
// The parsed map is then set to the internal spec map.
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

	m := make(models.SpecMap)
	if err = unmarshaller(data, &m); err != nil {
		return
	}

	if p.m == nil {
		p.m = NewSafeSpecMap(m)
	} else {
		p.m.Update(m)
	}

	return
}

// Spec returns the internal spec map instance.
func (p *baseProvider) Spec() *SafeSpecMap {
	return p.m
}
