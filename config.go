package main

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	DBPath string `yaml:"db_path"`
	TLS TLS `yaml:"tls"`
}

type TLS struct {
	CertPath string `yaml:"cert_path"`
	KeyPath string `yaml:"key_path"`
	LetsEncrypt bool `yaml:"letsencrypt"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
