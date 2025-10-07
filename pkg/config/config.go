package config

import (
	"log"

	"github.com/spf13/viper"
)

var Env *viper.Viper

func Load() {
	Env = viper.New()
	Env.SetConfigType("yaml")
	Env.SetConfigName("env")
	Env.AddConfigPath("../config/")
	Env.AddConfigPath("pkg/config")
	if err := Env.ReadInConfig(); err != nil {
		log.Println("Error while parsing config file:", err.Error())
	}
}

func GetConfig() *viper.Viper {
	return Env
}
