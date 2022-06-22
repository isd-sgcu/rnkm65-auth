package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSL      string `mapstructure:"ssl"`
}

type App struct {
	Port  int  `mapstructure:"port"`
	Debug bool `mapstructure:"debug"`
}

type ChulaSSO struct {
	Host         string `mapstructure:"host"`
	DeeAppID     string `mapstructure:"app-id"`
	DeeAppSecret string `mapstructure:"app-secret"`
}

type Jwt struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int32  `mapstructure:"expires_in"`
	Issuer    string `mapstructure:"issuer"`
}

type Config struct {
	Database Database `mapstructure:"database"`
	App      App      `mapstructure:"app"`
	ChulaSSO ChulaSSO `mapstructure:"chula-sso"`
	Jwt      Jwt      `mapstructure:"jwt"`
}

func LoadConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error occurs while reading the config")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "error occurs while unmarshal the config")
	}

	return
}
