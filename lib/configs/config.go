package configs

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type OutlineAPIConfig struct {
	API_URL string
}

func LoadConfigs() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	godotenv.Load(".env")

	return nil
}

func NewOutlineConfigs() OutlineAPIConfig {
	outline_cnf := OutlineAPIConfig{
		API_URL: os.Getenv("API_URL"),
	}
	return outline_cnf
}

func NewDbConfigs() DatabaseConfig {
	cnf := DatabaseConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.username"),
		Name:     viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	return cnf
}
