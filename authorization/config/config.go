package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App    `json:"app"`
		Log    `json:"logger"`
		JWT    `json:"jwt"`
		Server `json:"server"`
		UserDB `json:"user_db"`
	}
	App struct {
		AppEnviroment string `env-required:"true" json:"enviroment" env:"ENVIROMENT"`
	}
	Log struct {
		Filepath string `env-required:"false" json:"logfile_path" env:"LOGFILE_PATH"`
	}
	JWT struct {
		Secret     string        `env-required:"true" env:"JWT_SECRET"`
		AccessTTL  time.Duration `env-required:"true" json:"access_ttl" env:"ACCESS_TTL"`
		RefreshTTL time.Duration `env-required:"true" json:"refresh_ttl" env:"REFRESH_TTL"`
	}
	Server struct {
		Port             string        `env-required:"true" json:"port" env:"SERVER_PORT"`
		Host             string        `env-required:"true" json:"host" env:"SERVER_HOST"`
		ShutdownTimeout  time.Duration `env-required:"true" json:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
		OperationTimeout time.Duration `env-required:"true" json:"operation_timeout" env:"OPERATION_TIMEOUT"`
	}
	UserDB struct {
		DsnString string `env-required:"true" env:"DSN_PG_USER_STRING"`
	}
)

func MustLoadConfig(configPath string) *Config {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		log.Fatalf("Can not parse config file: %s\n", err.Error())
	}
	if err := cleanenv.UpdateEnv(cfg); err != nil {
		log.Fatalf("Can not update env: %s\n", err.Error())
	}
	return cfg
}
