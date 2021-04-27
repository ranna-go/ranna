package namespace

import "github.com/ranna-go/ranna/pkg/random"

type RandomProvider struct{}

func NewRandomProvider() *RandomProvider {
	return &RandomProvider{}
}

func (p *RandomProvider) Get() (string, error) {
	return random.GetRandBase64Str(32)
}
