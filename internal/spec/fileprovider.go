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
	m        models.SpecMap
}

func NewFileProvider(fileName string) *FileProvider {
	return &FileProvider{fileName: fileName, m: nil}
}

func (fp *FileProvider) Load() (err error) {
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

	err = unmarshaller(data, &fp.m)
	return
}

func (fp *FileProvider) Spec() models.SpecMap {
	return fp.m
}
