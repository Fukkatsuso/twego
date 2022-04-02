package main

import (
	"github.com/Fukkatsuso/twego/cmd"
)

func main() {
	cmd.Execute()

	// bearerToken := os.Getenv("TWITTER_BEARER_TOKEN")

	// rules := []SearchRule{
	// 	{
	// 		Value: "golang -is:retweet",
	// 		Tag:   "golang",
	// 	},
	// }

	// res, err := AddSearchRules(bearerToken, rules)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("add: %+v\n", res)

	// ids := []string{"1509899609703071747"}
	// err := DeleteSearchRules(bearerToken, ids)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// res, err := GetSearchRules(bearerToken)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Printf("rules: %+v\n", res)

	// done := make(chan struct{})
	// defer close(done)

	// stream := GetTweetStream(done, bearerToken)

	// w := tabwriter.NewWriter(os.Stdout, 0, 2, 0, ' ', 0)
	// for {
	// 	select {
	// 	case tweet, ok := <-stream:
	// 		if !ok {
	// 			fmt.Println("stream is closed")
	// 			return
	// 		}

	// 		now := time.Now().Format("2006/01/02 15:04:05")
	// 		texts := strings.Split(tweet.Data.Text, "\n")
	// 		for i, text := range texts {
	// 			if i == 0 {
	// 				fmt.Fprintln(w, now, "\t", text)
	// 			} else {
	// 				fmt.Fprintln(w, "\t", text)
	// 			}
	// 		}
	// 		w.Flush()
	// 	}
	// }
}
