package internal

import (
	"github.com/spf13/viper"
)

func InitViper() {
	// allow `GUMDROP_*` environment variables
	viper.SetEnvPrefix("gumdrop")

	// begin default configurations
	viper.SetDefault("BaseDir", ".")
	viper.BindEnv("BaseDir")

	viper.SetDefault("Address", ":8080")
	viper.BindEnv("Address")

	viper.SetDefault("FileMode", 0644)
	viper.BindEnv("FileMode")

	viper.SetDefault("PathMode", 0755)
	viper.BindEnv("PathMode")

	viper.SetDefault("LogFilePath", "")
	viper.BindEnv("LogFilePath")

	viper.SetDefault("Tokens", []string{"superSecretToken", "someOtherEquallySuperSecretToken"})
	// end default configurations

	// define config file parameters and locations to search
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/gumdrop")
	viper.AddConfigPath("$HOME/.gumdrop")
	viper.AddConfigPath("$HOME")

	// process/read the config
	_ = viper.ReadInConfig() // Find and read the config file
}
