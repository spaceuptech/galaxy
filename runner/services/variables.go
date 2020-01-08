package services

// Config is used to configure provider
type Config struct {
	Providers []string

	// For Digital Ocean
	DOToken string
	Region  string
}

// TypeProvider is the provider/cloud vendor type
type TypeProvider string

// ProviderDO is a constant of type digitalOcean i.e. "do"
const (
	ProviderDO TypeProvider = "do"
)
