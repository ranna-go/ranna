package spec

type Provider interface {
	Load() error
	Spec() *SafeSpecMap
}
