package spec

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/ranna/internal/models"
)

type FileProvider struct {
	fileName string
}

func NewFileProvider(fileName string) *FileProvider {
	return &FileProvider{fileName}
}

func (fp *FileProvider) Load() (m models.SpecMap, err error) {
	var unmarshaller func([]byte, interface{}) error

	switch strings.ToLower(path.Ext(fp.fileName)) {
	case ".yml", ".yaml":
		unmarshaller = yaml.Unmarshal
	case ".json":
		unmarshaller = json.Unmarshal
	default:
		err = errors.New("unsupported file type")
		return
	}

	data, err := os.ReadFile(fp.fileName)
	if err != nil {
		return
	}

	// m = make(map[string]*sandbox.Spec)
	err = unmarshaller(data, &m)
	return
}
