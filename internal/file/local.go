package file

type DummyFileProvider struct{}

func NewDummyFileProvider() *DummyFileProvider {
	return &DummyFileProvider{}
}

func (lf *DummyFileProvider) CreateDirectory(path string) error {
	return nil
}

func (lf *DummyFileProvider) CreateFileWithContent(path, content string) error {
	return nil
}

func (lf *DummyFileProvider) DeleteDirectory(path string) error {
	return nil
}
