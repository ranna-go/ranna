package namespace

type DummyProvider struct {
	staticName string
}

func NewDummyProvider(staticName string) *DummyProvider {
	return &DummyProvider{staticName}
}

func (p *DummyProvider) Get() (string, error) {
	return p.staticName, nil
}
