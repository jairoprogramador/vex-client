package vos

const (
	DefaultProvider = "local"
)

type Provider struct {
	Name string
}

func NewProvider(name string) *Provider {
	if name == "" {
		name = DefaultProvider
	}
	return &Provider{Name: name}
}