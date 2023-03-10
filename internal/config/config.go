package config

import (
	"path/filepath"
	"runtime"
	"sync"

	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	URL_VERSION = "/sys/version"
	URL_OPENAPI = "/openapi.json"
	URL_API_DOC = "/api"
)

type Config struct {
	Listen string `mapstructure:"listen" validate:"required"`
}

var (
	c    *Config
	once sync.Once
)

func Get() *Config {
	once.Do(setup)
	return c
}

func setup() {
	// get project root path
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")

	viper.SetEnvPrefix("WGS")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(projectRoot + "/configs")
	viper.AddConfigPath("/etc/wgs/")
	viper.AddConfigPath(projectRoot)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("read config failed")
	}

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal().Err(err).Msg("unmarshal config failed")
	}

	v := validator.New()
	if err := v.Struct(c); err != nil {
		log.Fatal().Err(err).Msg("validate config failed")
	}
	log.Info().Msg("load configs success")
}
