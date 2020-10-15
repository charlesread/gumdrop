package internal

import (
	"github.com/spf13/viper"
)

func validateToken(token string) bool {
	for _, v := range viper.GetStringSlice("Tokens") {
		if v == token {
			return true
		}
	}
	return false
}
