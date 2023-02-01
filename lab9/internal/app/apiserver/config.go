package apiserver

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{}
}
