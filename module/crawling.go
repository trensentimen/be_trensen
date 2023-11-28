package betrens

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dghubble/oauth1"
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
	// bearerToken := "AAAAAAAAAAAAAAAAAAAAALH7rAEAAAAAaclHSWIADYKfl6W6QRqP%2BH5rS90%3D2GCLNKms2djl39sCNi06I31GetaeuVFqSXs32Lqj7VShMrO6fH"
	accessToken := "1727210544716029952-e8kPGx6M1LS7Dv1rVaedY7Tc2oBNpw"
	accessTokenSecret := "ivqudihyOyPjgmJ9nqOVSelDN8iJhnJBzns73bMLH0aKw"
	// Set the API endpoint and parameters
	apiURL := "https://api.twitter.com/1.1/search/tweets.json?q=nasa"

	// Create OAuth1 client
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Make the HTTP request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response status and body
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:")
	fmt.Println(string(body))
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
