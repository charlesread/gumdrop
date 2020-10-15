package internal

import (
	"github.com/spf13/viper"
)

func setProcessResult(
	pr *processResult,
	err error,
	msg string,
	success bool,
	statusCode int) {
	pr.err = err
	pr.msg = msg
	pr.success = success
	pr.statusCode = statusCode
}

func validateToken(token string) bool {
	for _, v := range viper.GetStringSlice("Tokens") {
		if v == token {
			return true
		}
	}
	return false
}
