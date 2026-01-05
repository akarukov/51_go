package config

type Config struct {
	ServerAddr string
}

func GetConfig() *Config {
	return &Config{
		ServerAddr: "localhost:8080",
	}
}
