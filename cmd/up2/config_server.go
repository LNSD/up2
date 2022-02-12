package main

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	// Server host. Default: ""
	Host string

	// Server port. Default: 8080
	Port int `validate:"min=1,max=65535"`
}

func (sc ServerConfig) newConfigParser() *viper.Viper {
	parser := newConfigParser()

	parser.BindEnv("host")
	parser.SetDefault("host", "")

	parser.BindEnv("port")
	parser.SetDefault("port", 8080)

	return parser
}

func (sc *ServerConfig) LoadConfig() error {
	parser := sc.newConfigParser()
	if err := parser.Unmarshal(sc); err != nil {
		return err
	}
	return nil
}

func (sc ServerConfig) ValidateConfig() error {
	if err := validateConfig(sc); err != nil {
		return err
	}
	return nil
}

func (sc *ServerConfig) LoadAndValidateConfig() error {
	if err := sc.LoadConfig(); err != nil {
		return err
	}

	if err := sc.ValidateConfig(); err != nil {
		return err
	}

	return nil
}
