package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server `json:"server"`
		JWT    `json:"jwt"`
		Log    `json:"log"`
		PostDB `json:"post-db"`
	}
	Server struct {
		Host             string        `env-requires:"true" json:"host" env:"APP_HOST"`
		Port             string        `env-requires:"true" json:"port" env:"APP_PORT"`
		ShutdownTimeout  time.Duration `env-requires:"true" json:"shutdown-timeout" env:"JWT_SECRET_KEY"`
		OperationTimeout time.Duration `env-requires:"true" json:"operation-timeout" env:"JWT_SECRET_KEY"`
	}
	Log struct {
		Logpath string `env-requires:"true" json:"logfile-path" env:"LOGFILE_PATH"`
		Env     string `env-requires:"true" json:"env" env:"ENV"`
	}
	JWT struct {
		Secret    string `env-requires:"true" env:"JWT_SECRET_KEY"`
		AuthRoute string `env-requires:"true" json:"auth-route" env:"AUTH_ROUTE"`
	}
	PostDB struct {
		Dsn string `env-requires:"true" env:"DSN_PG_POST_DB"`
	}
)

func MustLoadConfig(configPath string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatal("Can't parse config: ", err.Error())
	}
	if err := cleanenv.UpdateEnv(cfg); err != nil {
		log.Fatal("Can't update enviroment: ", err.Error())
	}
	return cfg
}
