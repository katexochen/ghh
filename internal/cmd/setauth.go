package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	tokenEnvVar    = "GHH_TOKEN"
	configFilePath = "ghh/settings.json"
)

// NewSetAuthCmd creates a new command for setting the GitHub personal access token.
func NewSetAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-auth",
		Short: "Set the GitHub personal access token",
		RunE:  setupAuth,
	}
	return cmd
}

func setupAuth(_ *cobra.Command, _ []string) error {
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

	settings := settings{Token: token}

	file, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Join(configDir, "ghh"), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(configDir, configFilePath), file, 0o600); err != nil {
		return err
	}

	fmt.Println("Successfully saved token.")
	return nil
}

func readTokenFromUserInput() (string, error) {
	fmt.Println("Please enter your GitHub personal access token or restart -set-auth with the GHH_TOKEN environment variable set.")
	fmt.Print("Token: ")
	byteToken, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}

	return string(byteToken), nil
}

type settings struct {
	Token string `json:"token"`
}

func getToken() (string, error) {
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

	var settings settings
	if err := json.Unmarshal(file, &settings); err != nil {
		return "", err
	}

	return settings.Token, nil
}
