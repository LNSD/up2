package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"strings"
)

func newConfigParser() *viper.Viper {
	parser := viper.New()

	parser.SetEnvPrefix("UP2")
	parser.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	parser.AllowEmptyEnv(true)
	parser.AutomaticEnv()

	return parser
}

func validateConfig(config interface{}) error {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return err
	}
	return nil
}
