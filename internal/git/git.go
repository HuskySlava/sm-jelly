package git

type Config struct {
	RunDir string
}

type Git struct {
	Config *Config
}

func New(config *Config) *Git {
	return &Git{
		Config: config,
	}
}
