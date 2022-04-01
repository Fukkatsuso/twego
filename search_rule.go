package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type SearchRule struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

func AddSearchRules(bearerToken string, rules []SearchRule) ([]SearchRule, error) {
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

	// convert the response to []SearchRule
	var res struct {
		Data []SearchRule `json:"data"`
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
