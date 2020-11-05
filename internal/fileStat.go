package internal

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func getStat(directory string, file string) (os.FileInfo, error) {

	f, err := os.Open(filepath.Join(viper.GetString("BaseDir"), directory, file))
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	Log.Printf("stat: %v\n", stat)

	return stat, nil

}
