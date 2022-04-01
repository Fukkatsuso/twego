package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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

func GetTweetStream(bearerToken string) <-chan Tweet {
	stream := make(chan Tweet)

	go func() {
		defer close(stream)

		const endpoint = "https://api.twitter.com/2/tweets/search/stream"

		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)

		for {
			var tweet Tweet
			if err := decoder.Decode(&tweet); err != nil {
				fmt.Println(err)
				return
			}

			stream <- tweet
		}
	}()

	return stream
}
