package conf

import (
	"log"

	"github.com/spf13/viper"
)

// LoadEnvFromFile ...
func LoadEnvFromFile(file string) error {
	log.Printf("using config file at path %s \n", file)
	viper.SetConfigFile(file)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// Get ...
func Get(key string) string {
	return viper.GetString(key)
}

// Set ...
func Set(key string, value interface{}) {
	viper.Set(key, value)
}
