package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port               string `mapstructure:"PORT"`
	DBUrl              string `mapstructure:"DB_URL"`
	ServiceAccount     string `mapstructure:"SERVICE_ACCOUNT_PATH"`
	BucketName         string `mapstructure:"BUCKET_NAME"`
	TokenPath          string `mapstructure:"TOKEN_PATH"`
	FolderId           string `mapstructure:"FOLDER_ID"`
	ClientId           string `mapstructure:"CLIENT_ID"`
	ClientSecret       string `mapstructure:"CLIENT_SECRET"`
	GdriveRefreshToken string `mapstructure:"GDRIVE_API_REFRESH_TOKEN"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./config/env")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}

	return
}
