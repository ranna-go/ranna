package spec

import (
	"os"
	"path"
	"strings"
)

// FileProvider implements Provider retrieving a
// spec definitions from a file on the FS.
type FileProvider struct {
	*baseProvider
	fileName string
}

// NewFileProvider returns a new FileProvider reading
// the given fileName.
func NewFileProvider(fileName string) *FileProvider {
	return &FileProvider{baseProvider: newBaseProvider(), fileName: fileName}
}

func (t *FileProvider) Load() (err error) {
	data, err := os.ReadFile(t.fileName)
	if err != nil {
		return err
	}

	return t.parseAndSet(data, strings.ToLower(path.Ext(t.fileName)))
}
