package file

import "os"

type LocalFileProvider struct{}

func NewLocalFileProvider() *LocalFileProvider {
	return &LocalFileProvider{}
}

func (t *LocalFileProvider) CreateDirectory(path string) error {
	return os.MkdirAll(path, os.ModeDir|os.ModePerm)
}

func (t *LocalFileProvider) CreateFileWithContent(path, content string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func (t *LocalFileProvider) DeleteDirectory(path string) error {
	return os.RemoveAll(path)
}
