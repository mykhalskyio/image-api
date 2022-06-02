package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres struct {
		Port    int    `yaml:"port"`
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Name    string `yaml:"name"`
		Sslmode string `yaml:"sslmode"`
	} `yaml:"postgres"`

	RabbitMQ struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"rabbitmq"`
}

var (
	once sync.Once
	cfg  *Config
)

// get config
func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		if err := cleanenv.ReadConfig("config.yml", cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			log.Fatalln(err, help)
		}

	})
	return cfg
}
