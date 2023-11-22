package betrens

import (
	"context"
	"fmt"

	twitterscraper "github.com/n0madic/twitter-scraper"
)

func CrawlingTweet() {
	scraper := twitterscraper.New()
	err := scraper.Login("TrenSentimen", "@TrenS071023")
	if err != nil {
		panic(err)
	}

	fmt.Println(scraper.IsLoggedIn())
	for tweet := range scraper.SearchTweets(context.Background(),
		"twitter scraper data -filter:retweets", 50) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		fmt.Println(tweet.Text)
	}

}
