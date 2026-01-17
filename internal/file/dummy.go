package file

type DummyFileProvider struct{}

func NewDummyFileProvider() *DummyFileProvider {
	return &DummyFileProvider{}
}

func (t *DummyFileProvider) CreateDirectory(path string) error {
	return nil
}

func (t *DummyFileProvider) CreateFileWithContent(path, content string) error {
	return nil
}

func (t *DummyFileProvider) DeleteDirectory(path string) error {
	return nil
}
