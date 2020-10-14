package internal

import (
	"fmt"
	"github.com/spf13/viper"
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

	// process/read the config, panic if it doesn't work
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
