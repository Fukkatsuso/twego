package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rulesAddCmd.Flags().StringVarP(&rulesAddFlags.Tag, "tag", "t", "", "Tag the rule added")

	rulesCmd.AddCommand(
		rulesAddCmd,
		rulesDeleteCmd,
		rulesListCmd,
	)

	rootCmd.AddCommand(rulesCmd)
}

var rulesCmd = &cobra.Command{
	Use: "rules",
}

var rulesAddFlags struct {
	Tag string
}

var rulesAddCmd = &cobra.Command{
	Use:  "add",
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAuth("rules add")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		bearerToken := config.BearerToken

		rules := []Rule{
			{
				Value: args[0],
				Tag:   rulesAddFlags.Tag,
			},
		}

		res, err := AddRules(bearerToken, rules)
		if err != nil {
			return err
		}

		printRules(res)

		return nil
	},
}

var rulesDeleteCmd = &cobra.Command{
	Use:  "delete",
	Args: cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAuth("rules delete")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		bearerToken := config.BearerToken

		ids := args

		err := DeleteRules(bearerToken, ids)
		if err != nil {
			return err
		}

		for _, id := range ids {
			fmt.Println(id)
		}

		return nil
	},
}

var rulesListCmd = &cobra.Command{
	Use: "list",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAuth("rules list")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		bearerToken := config.BearerToken

		rules, err := GetRules(bearerToken)
		if err != nil {
			return err
		}

		printRules(rules)

		return nil
	},
}

func printRules(rules []Rule) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "ID\tVALUE\tTAG\n")
	for _, rule := range rules {
		fmt.Fprintf(w, "%s\t%s\t%s\n", rule.ID, rule.Value, rule.Tag)
	}
}

type Rule struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

func AddRules(bearerToken string, rules []Rule) ([]Rule, error) {
	const endpoint = "https://api.twitter.com/2/tweets/search/stream/rules"

	// convert rules to request body
	type addElem struct {
		Value string `json:"value"`
		Tag   string `json:"tag,omitempty"`
	}
	var add struct {
		Add []addElem `json:"add"`
	}
	for _, rule := range rules {
		elem := addElem{
			Value: rule.Value,
			Tag:   rule.Tag,
		}
		add.Add = append(add.Add, elem)
	}

	js, err := json.Marshal(add)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewBuffer(js)

	// create new request
	req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Add("Content-type", "application/json")

	// send the request
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// error handling
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%s", string(body))
	}

	// convert the response to []Rule
	var res struct {
		Data []Rule `json:"data"`
		Meta struct {
			Sent    time.Time `json:"sent"`
			Summary struct {
				Created    int `json:"created"`
				NotCreated int `json:"not_created"`
				Valid      int `json:"valid"`
				Invalid    int `json:"invalid"`
			} `json:"summary"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func DeleteRules(bearerToken string, ids []string) error {
	const endpoint = "https://api.twitter.com/2/tweets/search/stream/rules"

	// convert ids to request body
	var delete struct {
		Delete struct {
			Ids []string `json:"ids"`
		} `json:"delete"`
	}
	delete.Delete.Ids = ids

	js, err := json.Marshal(delete)
	if err != nil {
		return err
	}

	reqBody := bytes.NewBuffer(js)

	// create new request
	req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Add("Content-type", "application/json")

	// send the request
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// error handling
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}

	// convert the response
	var res struct {
		Meta struct {
			Sent    time.Time `json:"sent"`
			Summary struct {
				Deleted    int `json:"deleted"`
				NotDeleted int `json:"not_deleted"`
			} `json:"summary"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}

	return nil
}

func GetRules(bearerToken string) ([]Rule, error) {
	const endpoint = "https://api.twitter.com/2/tweets/search/stream/rules"

	// create new request
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// send the request
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// error handling
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", string(body))
	}

	// convert the response to []Rule
	var res struct {
		Data []Rule `json:"data"`
		Meta struct {
			Sent        time.Time `json:"sent"`
			ResultCount int       `json:"result_count"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res.Data, nil
}
