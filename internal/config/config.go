package config

import (
 "fmt"
 "os"

 "gopkg.in/yaml.v3"
)

type Config struct {
 Workers int `yaml:"workers"`
 WatchPath string `yaml:"watch_path"`
 Playlist string `yaml:"playlist"`
 R2 R2Config `yaml:"r2"`
}

type R2Config struct {
 Endpoint string `yaml:"endpoint"`
 Bucket string `yaml:"bucket"`
 AccessKey string `yaml:"access_key"`
 SecretKey string `yaml:"secret_key"`
 PublicURL string `yaml:"public_url"`
}

func Load(path string) (*Config, error) {
 data, err := os.ReadFile(path)
 if err != nil {
  return nil, fmt.Errorf("failed to read config: %w", err)
 }

 var cfg Config

 if err := yaml.Unmarshal(data, &cfg); err != nil {
  return nil, fmt.Errorf("failed to parse config: %w", err)
 }

 if cfg.Workers <= 0 {
  cfg.Workers = 64
 }

 return &cfg, nil
}