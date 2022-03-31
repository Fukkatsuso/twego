package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type SearchRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

func AddSearchRules(bearerToken string, rules []SearchRule) error {
	const endpoint = "https://api.twitter.com/2/tweets/search/stream/rules"

	js, err := json.Marshal(struct {
		Add []SearchRule `json:"add"`
	}{
		Add: rules,
	})
	if err != nil {
		return err
	}

	reqBody := bytes.NewBuffer(js)

	req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Add("Content-type", "application/json")

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

	var res struct {
		Data []struct {
			Value string `json:"value"`
			Tag   string `json:"tag,omitempty"`
			ID    string `json:"id"`
		} `json:"data"`
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
		return err
	}
	fmt.Printf("res: %+v\n", res)

	return nil
}

func DeleteSearchRules(bearerToken string, ids []string) error {
	const endpoint = "https://api.twitter.com/2/tweets/search/stream/rules"

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

	req, err := http.NewRequest(http.MethodPost, endpoint, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Add("Content-type", "application/json")

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
	fmt.Printf("res: %+v\n", res)

	return nil
}

func GetSearchRules(bearerToken string) error {
	const endpoint = "https://api.twitter.com/2/tweets/search/stream/rules"

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

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

	var res struct {
		Data []struct {
			ID    string `json:"id"`
			Value string `json:"value"`
			Tag   string `json:"tag,omitempty"`
		} `json:"data"`
		Meta struct {
			Sent        time.Time `json:"sent"`
			ResultCount int       `json:"result_count"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}
	fmt.Printf("res: %+v\n", res)

	return nil
}

func main() {
	bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	// rules := []SearchRule{
	// 	{
	// 		Value: "golang -is:retweet",
	// 		Tag:   "golang",
	// 	},
	// }

	// err := AddSearchRules(bearerToken, rules)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// ids := []string{"1509544617284280327"}
	// err := DeleteSearchRules(bearerToken, ids)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	err := GetSearchRules(bearerToken)
	if err != nil {
		fmt.Println(err)
	}
}
