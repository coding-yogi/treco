/*
Package conf to handle configuration files
*/
package conf

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config ...
type Config struct {
	Build        string
	Environment  string
	Jira         string
	ReportFile   string
	ReportFormat string
	Service      string
	TestType     string
	Coverage     string
}

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
	if viper.GetString(key) != "" {
		return viper.GetString(key)
	}

	return os.Getenv(key)
}

// Set ...
func Set(key string, value interface{}) {
	viper.Set(key, value)
}
