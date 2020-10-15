package internal

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var Log *log.Logger = nil

func InitLogger() {

	var w = os.Stdout

	logFilePath := viper.GetString("LogFilePath")
	if logFilePath != "" {
		f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			w = f
		}
	}

	Log = log.New(w, "", log.LstdFlags)

}
