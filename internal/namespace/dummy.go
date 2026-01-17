package namespace

type DummyProvider struct {
	staticName string
}

func NewDummyProvider(staticName string) *DummyProvider {
	return &DummyProvider{staticName}
}

func (t *DummyProvider) Get() (string, error) {
	return t.staticName, nil
}
