package betrens

import (
	"context"
	"fmt"

	twitterscraper "github.com/n0madic/twitter-scraper"
	model "github.com/trensentimen/be_trensen/model"
)

func CrawlingTweet(topic model.Topic) (dataTopic []model.DataTopics, err error) {
	scraper := twitterscraper.New()
	err = scraper.Login("TrenSentimen", "@TrenS071023")
	if err != nil {
		panic(err)
	}

	fmt.Println(scraper.IsLoggedIn())
	for tweet := range scraper.SearchTweets(context.Background(),
		topic.Source.Value, 5) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}
		// Create a new DataTopics instance with the tweet text
		data := model.DataTopics{
			Text:    tweet.Text,
			Date:    tweet.Timestamp,
			Source:  "twitter",
			TopicId: topic.ID,
			// You can add more fields as needed
		}
		dataTopic = append(dataTopic, data)

		// fmt.Println(tweet.Text)
		// fmt.Println(tweet.Timestamp)
	}
	return dataTopic, err
}
