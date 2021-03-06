package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(streamCmd)
}

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream Tweets in real-time that match the rules you added",
	Long:  "Stream Tweets in real-time that match the rules you added.\nIf you haven't added any rules, you will not receive any Tweets.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkAuth("stream")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		bearerToken := config.BearerToken

		done := make(chan struct{})
		defer close(done)

		stream := GetTweetStream(done, bearerToken)

		w := tabwriter.NewWriter(os.Stdout, 0, 2, 0, ' ', 0)
		for {
			select {
			case tweet, ok := <-stream:
				if !ok {
					return errors.New("stream is closed")
				}

				now := time.Now().Format("2006/01/02 15:04:05")
				texts := strings.Split(tweet.Data.Text, "\n")
				for i, text := range texts {
					if i == 0 {
						fmt.Fprintln(w, now, "\t", text)
					} else {
						fmt.Fprintln(w, "\t", text)
					}
				}
				w.Flush()
			}
		}
	},
}

type Tweet struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
	MatchingRules []struct {
		ID  string `json:"id"`
		Tag string `json:"tag"`
	} `json:"matching_rules"`
}

func GetTweetStream(done <-chan struct{}, bearerToken string) <-chan Tweet {
	stream := make(chan Tweet)

	go func() {
		defer close(stream)

		const endpoint = "https://api.twitter.com/2/tweets/search/stream"

		// create new request
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

		// send the request
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// error handling
		if resp.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			var res struct {
				Title  string `json:"title"`
				Detail string `json:"detail"`
				Type   string `json:"type"`
				Status int    `json:"status"`
			}
			if err := json.Unmarshal(body, &res); err != nil {
				fmt.Println(err)
				return
			}

			fmt.Printf("Error: %+v\n", res)
			return
		}

		// decode the response to Tweet
		decoder := json.NewDecoder(resp.Body)
		for {
			decoded := make(chan Tweet)
			go func() {
				defer close(decoded)
				var tweet Tweet
				if err := decoder.Decode(&tweet); err != nil {
					fmt.Println(err)
					return
				}
				decoded <- tweet
			}()

			select {
			case <-done:
				return
			case tweet, ok := <-decoded:
				if !ok {
					fmt.Println("cannot read the decoded response")
					return
				}
				select {
				case <-done:
					return
				case stream <- tweet:
				}
			}
		}
	}()

	return stream
}
