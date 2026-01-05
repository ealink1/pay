package config

import (
	"fmt"
	"os"
	"strconv"
)

type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	DSN      string
}

func LoadMySQLConfig(configPath string) (MySQLConfig, error) {
	var cfg MySQLConfig

	if configPath == "" {
		configPath = "config.yaml"
	}

	if fc, err := LoadConfig(configPath); err == nil {
		if fc.SQL.Host != "" || fc.SQL.DSN != "" || fc.SQL.Database != "" || fc.SQL.User != "" || fc.SQL.Port != 0 || fc.SQL.Password != "" {
			cfg = fc.SQL
		} else {
			cfg = fc.MySQL
		}
	}

	cfg = mergeMySQLEnvOverrides(cfg)
	cfg = applyMySQLDefaults(cfg)

	if cfg.DSN != "" {
		return cfg, nil
	}
	if cfg.Password == "" {
		return cfg, fmt.Errorf("mysql password is required (set sql.password in %s, or MYSQL_PASSWORD, or MYSQL_DSN)", configPath)
	}

	cfg.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	return cfg, nil
}

func LoadMySQLConfigFromEnv() (MySQLConfig, error) {
	return LoadMySQLConfig("")
}

func mergeMySQLEnvOverrides(cfg MySQLConfig) MySQLConfig {
	if v := os.Getenv("MYSQL_DSN"); v != "" {
		cfg.DSN = v
	}
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("MYSQL_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Port = n
		}
	}
	if v := os.Getenv("MYSQL_USER"); v != "" {
		cfg.User = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("MYSQL_DATABASE"); v != "" {
		cfg.Database = v
	}
	return cfg
}

func applyMySQLDefaults(cfg MySQLConfig) MySQLConfig {
	if cfg.Host == "" {
		cfg.Host = "114.132.245.76"
	}
	if cfg.Port == 0 {
		cfg.Port = 3306
	}
	if cfg.User == "" {
		cfg.User = "root"
	}
	if cfg.Database == "" {
		cfg.Database = "test_db"
	}
	return cfg
}

func getenv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getenvInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultValue
}
