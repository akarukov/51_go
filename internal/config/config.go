package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	ServerAddr    string `yaml:"ServerAddr"`
	ReferenceAddr string `yaml:"ReferenceAddr"`
}

func GetConfig() *Config {
	var cfg Config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	argServerAddr := flag.String("a", "", "The server address in the format of host:port")
	argReferenceAddr := flag.String("b", "", "The reference address in the format of host:port")

	flag.Parse()

	if *argServerAddr != "" {
		cfg.ServerAddr = *argServerAddr
	}
	if *argReferenceAddr != "" {
		cfg.ReferenceAddr = *argReferenceAddr
	}

	return &cfg
}
