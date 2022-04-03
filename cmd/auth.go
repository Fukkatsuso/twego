package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	authCmd.Flags().StringVar(&authFlags.BearerToken, "bearer", "", "Set Twitter Bearer Token")

	rootCmd.AddCommand(authCmd)
}

var authFlags struct {
	BearerToken string
}

var authCmd = &cobra.Command{
	Use: "auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		if authFlags.BearerToken != "" {
			return AuthByBearerToken(authFlags.BearerToken, defaultConfigFilepath())
		}

		return errors.New("Bearer Token is required")
	},
}

// check BearerToken
func checkAuth(cmd string) error {
	if config.BearerToken == "" {
		return fmt.Errorf(`Bearer Token is not set. To use "%s" command, please set the token by "auth" command`, cmd)
	}
	return nil
}

func AuthByBearerToken(bearerToken string, configFilepath string) error {
	config.BearerToken = bearerToken
	return writeConfig(config, configFilepath)
}
