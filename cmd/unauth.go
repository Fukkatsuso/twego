package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unauthCmd)
}

var unauthCmd = &cobra.Command{
	Use: "unauth",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAuth("unauth")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return Unauth(defaultConfigFilepath())
	},
}

func Unauth(configFilepath string) error {
	config.BearerToken = ""
	return writeConfig(config, configFilepath)
}
