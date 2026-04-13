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
	bindEnvVariables(Env)
}

func GetConfig() *viper.Viper {
	return Env
}

func bindEnvVariables(env *viper.Viper) {
	env.BindEnv("FILE_BASE_PATH", "FILEBASEPATH")
	env.BindEnv("SERVICE_HOST", "SERVICEHOST")
	env.BindEnv("JWT_SECRET", "JWT_SECRET")

	//Bind Postgres env variables
	env.BindEnv("DB_HOST", "DBHOST")
	env.BindEnv("DB_PORT", "DBPORT")
	env.BindEnv("DB_USER", "DBUSER")
	env.BindEnv("DB_NAME", "DBNAME")
	env.BindEnv("DB_PASSWORD", "DBPASSWORD")

	//Bind MongoDB env variables
	env.BindEnv("MongoDB_USER", "MONGOUSER")
	env.BindEnv("MongoDB_PASSWORD", "MONGOPASSWORD")
	env.BindEnv("MongoDB_HOST", "MONGOHOST")
	env.BindEnv("MongoDB_PORT", "MONGOPORT")
	env.BindEnv("MongoDB_DB_NAME", "MONGODBNAME")

	// Bind S3 env variables
	env.BindEnv("S3_ENDPOINT", "AWS_ENDPOINT_URL")
	env.BindEnv("S3_ACCESS_KEY", "AWS_ACCESS_KEY_ID")
	env.BindEnv("S3_SECRET_KEY", "AWS_SECRET_ACCESS_KEY")
	env.BindEnv("S3_BUCKET_NAME", "AWS_S3_BUCKET_NAME")
	env.BindEnv("S3_REGION", "AWS_DEFAULT_REGION")

	//Bind SMTP Server env variables
	env.BindEnv("SMTP_HOST", "SMTP_HOST")
	env.BindEnv("SMTP_PORT", "SMTP_PORT")
	env.BindEnv("SMTP_FROM", "SMTP_FROM")
	env.BindEnv("SMTP_USERNAME", "SMTP_USERNAME")
	env.BindEnv("SMTP_PASSWORD", "SMTP_PASSWORD")
	env.BindEnv("OTP_SECRET", "OTP_SECRET")
}
