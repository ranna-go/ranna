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

func (fp *FileProvider) Load() (err error) {
	data, err := os.ReadFile(fp.fileName)
	if err != nil {
		return
	}

	err = fp.parseAndSet(data, strings.ToLower(path.Ext(fp.fileName)))
	return
}
