package main

import (
	"github.com/charlesread/gumdrop/internal"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func init() {
	log.Printf("Starting `gumdrop`...")
}

func main() {

	log.Printf("Address: %q\n", viper.GetString("Address"))
	log.Printf("BaseDir: %v\n", viper.GetString("BaseDir"))

	server := &internal.Server{}
	log.Fatal(http.ListenAndServe(viper.GetString("Address"), server))

}
