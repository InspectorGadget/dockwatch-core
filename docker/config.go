package docker

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	DOCKER_HOST string
}

func GetConfig() (*Config, error) {
	var config Config

	file, err := os.OpenFile("config.json", os.O_RDONLY, 0644)
	if err != nil {
		// Fallback, get from Environment Variable
		dockerHost := os.Getenv("DOCKWATCH_DOCKER_HOST")
		if dockerHost != "" {
			return &Config{
				DOCKER_HOST: dockerHost,
			}, nil
		}

		return &Config{}, errors.New("unable to open config file")
	}

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return &Config{}, errors.New("failed to decode config file")
	}

	return &config, nil
}
