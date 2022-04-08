package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unauthCmd)
}

var unauthCmd = &cobra.Command{
	Use:   "unauth",
	Short: fmt.Sprintf("Delete your Twitter Bearer Token saved in %s", defaultConfigFilepath()),
	Long:  fmt.Sprintf("Delete your Twitter Bearer Token saved in %s.", defaultConfigFilepath()),
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
