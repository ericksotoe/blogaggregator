package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	// Username string `json: "current_user_name"`
}

// func getConfigFilePath() (string, error)

// func write(cfg Config) error

func Read() (Config, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	filePath = filePath + configFileName

	req, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	json.Unmarshal(req, &config)
	fmt.Printf("%s is the url read from the file at %s\n", config.DbUrl, filePath)

	return config, nil
}
