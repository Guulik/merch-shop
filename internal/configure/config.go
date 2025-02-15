package configure

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env"`
	Port     int           `yaml:"port"`
	Timeout  time.Duration `yaml:"timeout"`
	TokenTTL time.Duration `yaml:"tokenTTL"`
	Postgres Postgres
}

type Postgres struct {
	Host     string `yaml:"host"`
	SQLPort  int    `yaml:"SQLPort"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"DBName"`
	SslMode  string `yaml:"sslMode"`
}

func MustLoadConfig() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	const key = "CONFIG_PATH"

	if v := os.Getenv(key); v != "" {
		return v
	}

	//default local config
	return "./internal/configure/local.yaml"
}
