package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	authCmd.Flags().StringVar(&authFlags.BearerToken, "bearer", "", "Set Twitter Bearer Token")
	authCmd.Flags().StringVar(&authFlags.APIKey, "key", "", "Set Twitter API Key")
	authCmd.Flags().StringVar(&authFlags.APISecret, "secret", "", "Set Twitter API Secret")

	rootCmd.AddCommand(authCmd)
}

var authFlags struct {
	BearerToken string
	APIKey      string
	APISecret   string
}

var authCmd = &cobra.Command{
	Use: "auth",
	Example: strings.Join([]string{
		"auth --bearer $TWITTER_BEARER_TOKEN",
		"auth --key $TWITTER_API_KEY --secret $TWITTER_API_SECRET",
	}, "\n"),
	RunE: func(cmd *cobra.Command, args []string) error {
		if authFlags.BearerToken != "" {
			return AuthByBearerToken(authFlags.BearerToken, defaultConfigFilepath())
		}

		if authFlags.APIKey != "" && authFlags.APISecret != "" {
			return AuthByConsumerKeys(authFlags.APIKey, authFlags.APISecret, defaultConfigFilepath())
		}

		return errors.New("Bearer Token or Consumer Keys (API Key & API Secret) is required")
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

func AuthByConsumerKeys(apiKey, apiSecret string, configFilepath string) error {
	bearerToken, err := getBearerTokenByBasicAuth(apiKey, apiSecret)
	if err != nil {
		return err
	}
	return AuthByBearerToken(bearerToken, configFilepath)
}

func getBearerTokenByBasicAuth(apiKey, apiSecret string) (string, error) {
	const endpoint = "https://api.twitter.com/oauth2/token"

	reqBody := bytes.NewBufferString("grant_type=client_credentials")

	// create new request
	req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	// send the request
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// error handling
	if resp.StatusCode != http.StatusOK {
		var res struct {
			Errors []struct {
				Code    int    `json:"code"`
				Label   string `json:"label,omitempty"`
				Message string `json:"message"`
			} `json:"errors"`
		}
		if err := json.Unmarshal(body, &res); err != nil {
			return "", err
		}

		return "", fmt.Errorf("Errors: %+v", res.Errors)
	}

	// convert the response to struct
	var res struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}
