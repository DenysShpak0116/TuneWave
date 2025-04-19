package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string       `yaml:"env" env-default:"local"`
	StoragePath string       `yaml:"storage_path" env-required:"true"`
	JwtSecret   string       `yaml:"jwt_secret" env-required:"true"`
	Http        HttpConfig   `yaml:"http"`
	Google      GoogleConfig `yaml:"google"`
	Mail        MailConfig   `yaml:"mail"`
	AWS         AWSConfig    `yaml:"aws"`
}

type AWSConfig struct {
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Bucket    string `yaml:"bucket"`
}

type MailConfig struct {
	StmpServer   string `yaml:"smtp_server"`
	SmtpPort     int    `yaml:"smtp_port"`
	FromMail     string `yaml:"from_email"`
	FromPassword string `yaml:"from_password"`
}

type HttpConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timout"`
}

type GoogleConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
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
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
