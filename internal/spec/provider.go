package spec

// Provider to load the spec map.
type Provider interface {

	// Load retrieves and parses a spec
	// map from given source.
	Load() error

	// Spec returns the loaded spec map.
	Spec() *SafeSpecMap
}
