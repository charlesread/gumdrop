package internal

import (
	"github.com/spf13/viper"
	"log"
)

func init() {
	// allow `GUMDROP_*` environment variables
	viper.SetEnvPrefix("gumdrop")

	// begin default configurations
	viper.SetDefault("BaseDir", ".")
	viper.BindEnv("BaseDir")

	viper.SetDefault("Address", ":8080")
	viper.BindEnv("Address")
	// end default configurations

	// define config file parameters and locations to search
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/gumdrop")
	viper.AddConfigPath("$HOME/.gumdrop")
	viper.AddConfigPath("$HOME")

	// process/read the config, panic if it doesn't work
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Printf("Config file not found, default config being used: %v\n", err)
	}
}
