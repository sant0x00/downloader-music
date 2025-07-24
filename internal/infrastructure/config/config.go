package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Download DownloadConfig `yaml:"download"`
	Scraping ScrapingConfig `yaml:"scraping"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type DownloadConfig struct {
	ConcurrentWorkers int    `yaml:"concurrent_workers"`
	RetryAttempts     int    `yaml:"retry_attempts"`
	TimeoutSeconds    int    `yaml:"timeout_seconds"`
	OutputDirectory   string `yaml:"output_directory"`
}

type ScrapingConfig struct {
	BaseURL              string        `yaml:"base_url"`
	DelayBetweenRequests time.Duration `yaml:"delay_between_requests"`
	UserAgent            string        `yaml:"user_agent"`
}

type LoggingConfig struct {
	Level      string `yaml:"level"`
	OutputFile string `yaml:"output_file"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{
		Download: DownloadConfig{
			ConcurrentWorkers: 8,
			RetryAttempts:     3,
			TimeoutSeconds:    30,
			OutputDirectory:   "~/Downloads/ClipesJW",
		},
		Scraping: ScrapingConfig{
			BaseURL:              "https://www.jw.org/pt/biblioteca/musica-canticos/clipes-musicais/",
			DelayBetweenRequests: time.Second,
			UserAgent:            "ClipesJW-Downloader/1.0",
		},
		Logging: LoggingConfig{
			Level:      "info",
			OutputFile: "downloader.log",
		},
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, SaveConfig(config, configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if config.Download.OutputDirectory[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		config.Download.OutputDirectory = filepath.Join(homeDir, config.Download.OutputDirectory[1:])
	}

	return config, nil
}

func SaveConfig(config *Config, configPath string) error {
	// Criar diretório se não existir
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
