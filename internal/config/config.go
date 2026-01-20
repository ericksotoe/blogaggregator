package config

import (
	"encoding/json"
	"os"
	"strings"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DbUrl    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

// retrieves the config filePath
func getConfigFilePath() (string, error) {
	filePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath = filePath + configFileName
	return filePath, nil
}

// converts the Config struct to json and writes it to our config file
func write(cfg Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	return nil
}

// converts the json found at filepath to the Config struct
func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	req, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	json.Unmarshal(req, &config)
	return config, nil
}

// sets the user at in the Config struct
func (c *Config) SetUser(name string) error {
	c.Username = name
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

// this function normalizes the text input and returns them
func cleanInput(text string) []string {
	text = strings.ToLower(text)
	words := strings.Fields(text)
	return words
}
