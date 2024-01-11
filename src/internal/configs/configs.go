package configs

import (
	"encoding/json"
	"flag"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	App           App
	TSM           TSM
	Commercetools Commercetools
}

type App struct {
	Name        string
	Version     string
	Environment string
}

type TSM struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Interval int
}

type Commercetools struct {
	API          string
	AUTH         string
	Import       string
	ProjectKey   string
	ClientID     string
	ClientSecret string
}

var cfg = &Config{}

func GetConfig() (*Config, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg, err := cfg.ReadConfig()
	return cfg, err
}

var path = flag.String("config", "./", "Location of the config file.")

func (c *Config) ReadConfig() (*Config, error) {
	viper.AddConfigPath(*path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	opt := viper.DecodeHook(
		jsonStringToStruct(c),
	)
	err = viper.Unmarshal(c, opt)
	return c, err
}

func jsonStringToStruct(m interface{}) func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
	return func(rf reflect.Kind, rt reflect.Kind, data interface{}) (interface{}, error) {
		if rf != reflect.String || rt != reflect.Struct {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return m, nil
		}

		err := json.Unmarshal([]byte(raw), &m)
		return m, err
	}
}
