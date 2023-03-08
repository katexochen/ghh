package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/term"
)

const (
	tokenEnvVar    = "GHH_TOKEN"
	configFilePath = "ghh/settings.json"
)

type Settings struct {
	Token string `json:"token"`
}

func GetToken() (string, error) {
	token := os.Getenv(tokenEnvVar)
	if token != "" {
		return token, nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	file, err := os.ReadFile(filepath.Join(configDir, configFilePath))
	if errors.Is(err, os.ErrNotExist) {
		return "", errors.New("no token found. Please set the GHH_TOKEN environment variable or run `ghh setauth`")
	} else if err != nil {
		return "", err
	}

	var settings Settings
	if err := json.Unmarshal(file, &settings); err != nil {
		return "", err
	}

	return settings.Token, nil
}

func SetupAuth() error {
	var token string

	token = os.Getenv(tokenEnvVar)
	if token == "" {
		var err error
		token, err = readTokenFromUserInput()
		if err != nil {
			return err
		}
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	settings := Settings{Token: token}

	file, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(configDir, "ghh"), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(configDir, configFilePath), file, 0o644); err != nil {
		return err
	}

	fmt.Println("Successfully saved token.")
	return nil
}

func readTokenFromUserInput() (string, error) {
	fmt.Println("Please enter your GitHub personal access token or restart auth with the GHH_TOKEN environment variable set.")
	fmt.Print("Token: ")
	byteToken, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	return string(byteToken), nil
}
