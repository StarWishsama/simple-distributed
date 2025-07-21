package util

import (
	"github.com/spf13/viper"
)

func InitViper(cfgName string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(cfgName)
	v.AddConfigPath(".")
	v.AddConfigPath("./.env")
	v.SetConfigType("toml")

	err := v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return v, nil
}
