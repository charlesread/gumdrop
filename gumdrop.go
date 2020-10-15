package main

import (
	"github.com/charlesread/gumdrop/internal"
	"github.com/spf13/viper"
	"net/http"
)

func main() {

	internal.InitViper()
	internal.InitLogger()

	internal.Log.Printf("Starting `gumdrop`...")
	internal.Log.Printf("Address: %q\n", viper.GetString("Address"))
	internal.Log.Printf("BaseDir: %v\n", viper.GetString("BaseDir"))
	internal.Log.Printf("LogFilePath: %v\n", viper.GetString("LogFilePath"))
	internal.Log.Printf("Tokens: %v\n", viper.GetStringSlice("Tokens"))

	server := &internal.Server{}
	internal.Log.Fatal(http.ListenAndServe(viper.GetString("Address"), server))

}
