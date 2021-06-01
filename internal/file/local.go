package file

import "os"

type LocalFileProvider struct{}

func NewLocalFileProvider() *LocalFileProvider {
	return &LocalFileProvider{}
}

func (lf *LocalFileProvider) CreateDirectory(path string) error {
	return os.MkdirAll(path, os.ModeDir)
}

func (lf *LocalFileProvider) CreateFileWithContent(path, content string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return
}

func (lf *LocalFileProvider) DeleteDirectory(path string) error {
	return os.RemoveAll(path)
}
