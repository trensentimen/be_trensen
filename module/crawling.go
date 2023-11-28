package betrens

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	twitterscraper "github.com/n0madic/twitter-scraper"
	model "github.com/trensentimen/be_trensen/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CrawlingTweet(topic model.Topic) (dataTopic []model.DataTopics, err error) {
	scraper := twitterscraper.New()
	f, _ := os.Open("cookies.json")
	// deserialize from JSON
	var cookies []*http.Cookie
	json.NewDecoder(f).Decode(&cookies)
	// load cookies
	scraper.SetCookies(cookies)
	// check login status
	fmt.Println(scraper.IsLoggedIn())

	// get data from dbvar doc model.Setting

	var docSetting model.Setting
	db := MongoConnect("MONGOSTRING", "trensentimen")
	_id, err := primitive.ObjectIDFromHex("656455b184d8d327072ba54b")
	if err != nil {
		fmt.Printf("Error get setting: %v", err)
	}
	docSetting.ID = _id
	setting, err := GetSetting(db, docSetting)

	for tweet := range scraper.SearchTweets(context.Background(),
		topic.Source.Value, setting.MaxTweeter) {
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
	}
	return dataTopic, err
}

func GetCookieTwitter() (err error) {
	scraper := twitterscraper.New()
	// err = scraper.Login("TrenSentimen", "@TrenS071023")
	if err = scraper.Login("TrenSentimen", "@TrenS071023", "trensentimen@gmail.com"); err != nil {
		fmt.Printf("Error logging in to Twitter: %v", err)
		return err
	}

	cookies := scraper.GetCookies()
	// serialize to JSON
	js, _ := json.Marshal(cookies)
	// save to file
	f, _ := os.Create("cookies.json")
	f.Write(js)
	return nil
}

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
		MaxResults(5). // Adjust the maximum results as needed
		Order("time")  // Specify the sorting order

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
			Source:  "youtube",
			TopicId: topic.ID,
			// You can add more fields as needed
		}
		dataTopic = append(dataTopic, data)
		// Process or store the extracted comment data
		fmt.Println(commentText) // Print the comment text for demonstration
	}

	return dataTopic, err
}
