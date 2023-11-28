package betrens

import (
	"context"
	"fmt"
	"time"

	anaconda "github.com/ChimeraCoder/anaconda"
	twitterscraper "github.com/n0madic/twitter-scraper"
	model "github.com/trensentimen/be_trensen/model"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CrawlingTweet2(topic model.Topic) (dataTopic []model.DataTopics, err error) {
	scraper := twitterscraper.New()
	// err = scraper.Login("TrenSentimen", "@TrenS071023")
	if err = scraper.Login("TrenSentimen", "@TrenS071023", "trensentimen@gmail.com"); err != nil {
		fmt.Printf("Error logging in to Twitter: %v", err)
		return nil, err
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
	scraper.Logout()
	return dataTopic, err
}

func CrawlingTweet(topic model.Topic) (dataTopic []model.DataTopics, err error) {
	consumerKey := "akA2uYm8PKzmB44f2NEQhMfkT"
	consumerSecret := "syWSZxb5dpIJVoIBj7yW7nc9xvzkN0nl3GJDmPCDKkt26qcP3f"
	accessToken := "1727210544716029952-e8kPGx6M1LS7Dv1rVaedY7Tc2oBNpw"
	accessSecret := "ivqudihyOyPjgmJ9nqOVSelDN8iJhnJBzns73bMLH0aKw"

	// Authenticate with Twitter API
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)

	// Set the search query parameters
	searchResult, _ := api.GetSearch("jokowi", nil)

	fmt.Print(searchResult.Metadata.CompletedIn)
	for _, tweet := range searchResult.Statuses {
		fmt.Print(tweet.Text)
	}

	return dataTopic, err
}

func CrawlingYoutube(topic model.Topic) (dataTopic []model.DataTopics, errM string, err error) {
	// AIzaSyC4yKKLe58P33lc_MlTelPmPbJkZcluT9Y
	apiKey := "AIzaSyC4yKKLe58P33lc_MlTelPmPbJkZcluT9Y"
	videoID := topic.Source.Value
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, "Apikey Bermasalah", err
		// panic(err)
	}

	call := youtubeService.CommentThreads.List([]string{"snippet"}).
		VideoId(videoID).
		MaxResults(5). // Adjust the maximum results as needed
		Order("time")  // Specify the sorting order

	response, err := call.Do()
	if err != nil {
		return nil, "id video youtube salah", err
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
			Source:  "youtube",
			TopicId: topic.ID,
			// You can add more fields as needed
		}
		dataTopic = append(dataTopic, data)
		// Process or store the extracted comment data
		fmt.Println(commentText) // Print the comment text for demonstration
	}

	return dataTopic, "", err
}
