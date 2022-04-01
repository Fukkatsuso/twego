package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

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

	// err := GetSearchRules(bearerToken)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 0, ' ', 0)

	stream := GetTweetStream(bearerToken)
	for {
		select {
		case tweet, ok := <-stream:
			if !ok {
				fmt.Println("stream is closed")
				return
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
}
