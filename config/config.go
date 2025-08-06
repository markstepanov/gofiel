package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadConfigFile() error {
	configFilename, err := getFilePath()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configFilename)

	if err != nil {
		return err
	}

	var config = Config{}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return err
	}

	err = verifyEssestialConfigValues(&config)

	if err != nil {
		return err
	}

	GlobalServerConfig = config
	return nil
}

func verifyEssestialConfigValues(config *Config) error {
	if config.Port == "" {
		return errors.New("port must be set in the config file")
	}

	if config.Basedir == "" {
		return errors.New("base-dir must be set in the config file")
	}

	if config.Adminusername == "" {
		return errors.New("admin-username must be set in the config file")
	}

	if config.Adminpassword == "" {
		return errors.New("admin-password must be set in the config file")
	}

	return nil
}

func getFilePath() (string, error) {
	args := os.Args
	argsLen := len(args) - 1
	for i, val := range args {
		if val == "-c" && i+1 <= argsLen {
			return args[i+1], nil
		}
	}
	return "", errors.New("config filename was not supplied (use gofiel -c <filename>)")

}
