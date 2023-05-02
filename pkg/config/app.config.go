package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	Port      string `mapstructure:"PORT"`
	DbHost    string `mapstructure:"DB_HOST"`
	DbPort    string `mapstructure:"DB_PORT"`
	DbUser    string `mapstructure:"DB_USER"`
	DbPass    string `mapstructure:"DB_PASS"`
	DbTLS     bool   `mapstructure:"DB_TLS"`
	DbName    string `mapstructure:"DB_NAME"`
	ApiSecret string `mapstructure:"API_SECRET"`
}

var Env string
var Config Configuration

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Configuration, err error) {
	env := GetEnv("ENV", "local")

	vp := viper.New()
	vp.AddConfigPath(".")
	vp.AddConfigPath("../")
	vp.AddConfigPath("../pkg/")
	vp.SetConfigName("config")

	err = vp.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("unable to get config for env:%s", env))
	}

	sub := vp.Sub(env)
	if sub == nil {
		panic(fmt.Sprintf("unable to get config for env:%s", env))
	}

	err = sub.Unmarshal(&config)
	if err != nil {
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}
	Env = env
	Config = config
	fmt.Printf("Configuration env:%s\n", env)
	return
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		if value == "" {
			return fallback
		}
		return value
	}
	return fallback
}
