package config

import (
	"errors"
	"log"
	"os"
)

func ReadConfigFile() (*Config, error) {
	configFilename, err := getFilePath()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(configFilename)

	return nil, nil
}

func getFilePath() (string, error) {
	args := os.Args
	argsLen := len(args) - 1
	for i := 0; i < argsLen; i++ {
		if args[i] == "-c" && i+1 <= argsLen {
			return args[i+1], nil
		}
	}
	return "", errors.New("config filename was not supplied (use gofiel -c <filename>)")

}
