package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	PostgresDsn    string        `yaml:"postgres_dsn" env-required:"true"`
	MongoDb_uri    string        `yaml:"mongoDb_uri"`
	Timeout        time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout    time.Duration `yaml:"iddle_timeout" env-default:"60s"`
	UserManager    int           `yaml:"userManager"`
	Notification   int           `yaml:"notification"`
	DataCollection int           `yaml:"dataCollection"`
	SearchManager  int           `yaml:"searchManager"`
	Review         int           `yaml:"review"`
	//HTTPserver  `yaml:"http_server"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}

//type HTTPserver struct {
//	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
//	IdleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
//	ports       `yaml:"ports"`
//}
//
//type ports struct {
//	UserManager    int `yaml:"userManager"`
//	Notification   int `yaml:"notification"`
//	DataCollection int `yaml:"dataCollection"`
//	SearchManager  int `yaml:"searchManager"`
//	Review         int `yaml:"review"`
//}
