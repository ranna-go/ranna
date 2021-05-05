package spec

import (
	"os"
	"path"
	"strings"
)

type FileProvider struct {
	*baseProvider
	fileName string
}

func NewFileProvider(fileName string) *FileProvider {
	return &FileProvider{baseProvider: newBaseProvider(), fileName: fileName}
}

func (fp *FileProvider) Load() (err error) {
	data, err := os.ReadFile(fp.fileName)
	if err != nil {
		return
	}

	err = fp.parseAndSet(data, strings.ToLower(path.Ext(fp.fileName)))
	return
}
