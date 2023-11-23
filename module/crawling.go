package betrens

import (
	"context"
	"fmt"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
	model "github.com/trensentimen/be_trensen/model"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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

func CrawlingYoutube(topic model.Topic) (dataTopic []model.DataTopics, err error) {
	// AIzaSyC4yKKLe58P33lc_MlTelPmPbJkZcluT9Y
	apiKey := "AIzaSyC4yKKLe58P33lc_MlTelPmPbJkZcluT9Y"
	videoID := topic.Source.Value
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	call := youtubeService.CommentThreads.List([]string{"snippet"}).
		VideoId(videoID).
		MaxResults(50). // Adjust the maximum results as needed
		Order("time")   // Specify the sorting order

	response, err := call.Do()
	if err != nil {
		panic(err)
	}

	for _, item := range response.Items {
		// Extract the top-level comment text
		commentText := item.Snippet.TopLevelComment.Snippet.TextDisplay
		publishedAtStr := item.Snippet.TopLevelComment.Snippet.PublishedAt
		publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
		if err != nil {
			panic(err)
		}
		unixTimestamp := publishedAt.Unix()
		data := model.DataTopics{
			Text:    commentText,
			Date:    unixTimestamp,
			Source:  "twitter",
			TopicId: topic.ID,
			// You can add more fields as needed
		}
		dataTopic = append(dataTopic, data)
		// Process or store the extracted comment data
		fmt.Println(commentText) // Print the comment text for demonstration
	}

	return dataTopic, err
}
