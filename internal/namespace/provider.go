package namespace

type Provider interface {
	Get() (string, error)
}
