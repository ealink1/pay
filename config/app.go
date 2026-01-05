package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

type LogConfig struct {
	Level            string        `yaml:"level"`
	Encoding         string        `yaml:"encoding"`
	OutputPaths      []string      `yaml:"output_paths"`
	ErrorOutputPaths []string      `yaml:"error_output_paths"`
	Development      bool          `yaml:"development"`
	File             LogFileConfig `yaml:"file"`
}

type LogFileConfig struct {
	Enabled     bool   `yaml:"enabled"`
	Dir         string `yaml:"dir"`
	Filename    string `yaml:"filename"`
	MaxAgeDays  int    `yaml:"max_age_days"`
	RotateHours int    `yaml:"rotate_hours"`
}

type TraceConfig struct {
	Enabled bool   `yaml:"enabled"`
	Header  string `yaml:"header"`
}

type AlipayConfig struct {
	NotifyURL string `yaml:"notify_url"`
	ReturnURL string `yaml:"return_url"`
}

type AlipayAppConfig struct {
	AppId           string `yaml:"appId"`
	PrivateKey      string `yaml:"privateKey"`
	AlipayPublicKey string `yaml:"alipayPublicKey"`
	NotifyURL       string `yaml:"notify_url"`
	ReturnURL       string `yaml:"return_url"`
}

type PayConfig struct {
	AlipaySandbox AlipayAppConfig `yaml:"alipaySandbox"`
	Alipay        AlipayAppConfig `yaml:"alipay"`
}

type AppConfig struct {
	SQL    MySQLConfig  `yaml:"sql"`
	Log    LogConfig    `yaml:"log"`
	Trace  TraceConfig  `yaml:"trace"`
	Pay    PayConfig    `yaml:"pay"`
	Alipay AlipayConfig `yaml:"alipay"`
	MySQL  MySQLConfig  `yaml:"mysql"`
}

func LoadConfig(configPath string) (AppConfig, error) {
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg := AppConfig{
		Log: LogConfig{
			Level:            "info",
			Encoding:         "json",
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			File: LogFileConfig{
				Enabled:     false,
				Dir:         "logs",
				Filename:    "app",
				MaxAgeDays:  14,
				RotateHours: 24,
			},
		},
		Trace: TraceConfig{
			Enabled: true,
			Header:  "X-Trace-Id",
		},
	}

	if data, err := os.ReadFile(filepath.Clean(configPath)); err == nil && len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return AppConfig{}, err
		}
	}

	if cfg.Trace.Header == "" {
		cfg.Trace.Header = "X-Trace-Id"
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Encoding == "" {
		cfg.Log.Encoding = "json"
	}
	if len(cfg.Log.OutputPaths) == 0 {
		cfg.Log.OutputPaths = []string{"stdout"}
	}
	if len(cfg.Log.ErrorOutputPaths) == 0 {
		cfg.Log.ErrorOutputPaths = []string{"stderr"}
	}
	if cfg.Log.File.Dir == "" {
		cfg.Log.File.Dir = "logs"
	}
	if cfg.Log.File.Filename == "" {
		cfg.Log.File.Filename = "app"
	}
	if cfg.Log.File.MaxAgeDays <= 0 {
		cfg.Log.File.MaxAgeDays = 14
	}
	if cfg.Log.File.RotateHours <= 0 {
		cfg.Log.File.RotateHours = 24
	}

	if cfg.Pay.AlipaySandbox.NotifyURL == "" {
		cfg.Pay.AlipaySandbox.NotifyURL = cfg.Alipay.NotifyURL
	}
	if cfg.Pay.AlipaySandbox.ReturnURL == "" {
		cfg.Pay.AlipaySandbox.ReturnURL = cfg.Alipay.ReturnURL
	}
	if cfg.Pay.Alipay.NotifyURL == "" {
		cfg.Pay.Alipay.NotifyURL = cfg.Alipay.NotifyURL
	}
	if cfg.Pay.Alipay.ReturnURL == "" {
		cfg.Pay.Alipay.ReturnURL = cfg.Alipay.ReturnURL
	}

	return cfg, nil
}
