package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	TokenTTL string   `yaml:"token_ttl"`
	DB       DbConfig `yaml:"db"`
}

type DbConfig struct {
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("no config path provided")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file does not exist at path " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
