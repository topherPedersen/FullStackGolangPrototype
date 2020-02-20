package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Initialze TwitterGo Client
func InitializeTwitterGoClient() (client *twittergo.Client) {

	// NOTE: Keys/Secrets/Tokens Below Have Been Scrubbed, Fake Placeholders Added

	// TophersTop100 App Key: SLDKFJSDFKLJSDJFLKJLKSJDFSDFLKJ
	// TophersTop100 App Secret: WEKRLJSDFSDFKSDJFDFKIIEEKKJSDFJKsdJFKSJDFKSDFsdfsdfsdfsdfsdf
	// My Personal Access Token (@topherPedersen): 982374982734-lskjdflksjdflkjasdfkjjhkjashdfkjhasdfhkj
	// My Personal Access Token Secret (@topherPedersen): lksdfoiuwdefjkhnskjdfhksjdhfsdfgsdhfjkj

	config := &oauth1a.ClientConfig{
		ConsumerKey:    "SLDKFJSDFKLJSDJFLKJLKSJDFSDFLKJ",
		ConsumerSecret: "WEKRLJSDFSDFKSDJFDFKIIEEKKJSDFJKsdJFKSJDFKSDFsdfsdfsdfsdfsdf",
	}
	userAccessToken := "982374982734-lskjdflksjdflkjasdfkjjhkjashdfkjhasdfhkj"
	userAccessTokenSecret := "lksdfoiuwdefjkhnskjdfhksjdhfsdfgsdhfjkj"
	user := oauth1a.NewAuthorizedConfig(userAccessToken, userAccessTokenSecret)
	client = twittergo.NewClient(config, user)
	return
}

// ServerResponse struct will be used for forming our JSON response
// which will be returned back to the client
type ServerResponse struct {
	Songs []string
}

func mainRoute(c echo.Context) error {

	// The serverResponse struct will be used to create our JSON
	// object which will be returned back to the client. This struct
	// will be populated with links to the top 100 songs returned from the
	// Twitter API call below.
	var serverResponse ServerResponse

	// Greet User
	fmt.Println("Welcome to Topher's Top 100")
	fmt.Println("App currently under development and coming soon...")
	fmt.Println("In the meantime, let's search twitter using keywords!")
	fmt.Println("")

	// Set Search Query
	searchQuery := "soundcloud.com"

	// Declare TwitterGo Variables
	var (
		err     error
		client  *twittergo.Client
		req     *http.Request
		resp    *twittergo.APIResponse
		results *twittergo.SearchResults
	)

	// Instantiate TwitterGo Client
	client = InitializeTwitterGoClient()

	// **************************************************************************
	// *** NOTE REGARDING tweet_mode=extended AND 140 VS 280 CHARACTER LIMITS ***
	// **************************************************************************
	//
	// When performing standard search queries without setting tweet_mode to extended,
	// the Twitter API will only return 140 character tweets. But setting tweet_mode
	// equal to extended will allow us to fetch full 280 character tweets. However,
	// retweets for some reason will still be truncated to 140 characters.
	// Also, when using tweet_mode=extended we need to use tweet.FullText instead of
	// tweet.Text. Furthermore, if you need to fetch the full 280 characters for
	// a retweet, fetch the original tweet as suggested by user dtabares in the
	// thread mentioned below.
	//
	// REFERENCE: "280 characters via REST API only works for full_text and tweet_mode=extended"
	// URL: https://github.com/sferik/twitter/issues/880
	//
	// **************************************************************************
	// **************************************************************************

	// **************************************************************************
	// TWITTER STANDARD 7 DAY SEARCH API DOCUMENTATION
	// **************************************************************************
	//
	// Important Information Regarding Search Query Parameters:
	//
	// https://developer.twitter.com/en/docs/tweets/search/api-reference/get-search-tweets
	//
	// **************************************************************************
	// **************************************************************************

	// Query Twitter API
	query := url.Values{}
	query.Set("q", searchQuery)
	// geocode=30,-97,250mi => 250 mile radius of Austin, TX
	url := fmt.Sprintf("/1.1/search/tweets.json?%v&tweet_mode=extended&count=100&geocode=30,-97,250mi&include_entities=true", query.Encode())
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err = client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	results = &twittergo.SearchResults{}
	err = resp.Parse(results)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}

	// Handle Query Results Returned from Twitter API
	for i, tweet := range results.Statuses() {

		// Get Tweet Info
		user := tweet.User()
		tweetText := tweet.FullText()
		twitterHandle := user.ScreenName()
		// REFERENCE (Working with Time in Golang): https://yourbasic.org/golang/day-month-year-from-time/
		tweetTime := tweet.CreatedAt()
		tweetYear, tweetMonth, tweetDay := tweetTime.Date()
		tweetUnixTimestamp := tweetTime.Unix()
		tweetId := tweet.Id()
		tweetIdStr := tweet.IdStr()
		fmt.Printf("%d) @%v: ", i+1, twitterHandle)
		fmt.Printf("%v\n", tweetText)
		fmt.Printf("Year: %d Month: %d Day: %d \n", tweetYear, tweetMonth, tweetDay)
		fmt.Printf("Timestamp: %d \n", tweetUnixTimestamp)
		fmt.Printf("TweetId: %d \n", tweetId)
		fmt.Printf("TweetIdStr: %v \n", tweetIdStr)

		// URLs
		URLMap := tweet.Entities().URLs()
		// We use a loop to print the extended URLS in case there are
		// no URLS to print. For example, if there are no attached URLs
		// to a tweet, attempting to print a URL will result in a panic
		// runtime error: index out of range.
		// However, by wrapping the print statements in a loop the
		// the program will not crash if the URLMap is empty.
		for _, URL := range URLMap {
			// fmt.Println(URL["expanded_url"])
			expandedUrlStr := fmt.Sprint(URL["expanded_url"])
			// songs = append(songs, expandedUrlStr)
			fmt.Println(expandedUrlStr)
			serverResponse.Songs = append(serverResponse.Songs, expandedUrlStr)
		}

		// Identify Retweets
		// ( This is done by determining if the Tweet contains the Prefix, "RT")
		isRetweet := strings.HasPrefix(tweetText, "RT")
		if isRetweet {
			fmt.Println("IS A RETWEET!!! IS A RETWEET!!! IS A RETWEET!!! IS A RETWEET!!! IS A RETWEET!!!")
		} else {
			fmt.Println("is not a retweet...")
		}

		// TODO: Identify Comments

		fmt.Println("----------------------------------------------------")
	}

	// Extra stuff...
	// callback := c.QueryParam("callback")

	// Now return the JSON back to the Client
	return c.JSON(http.StatusOK, &serverResponse)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", mainRoute)
	e.Logger.Fatal(e.Start(":1323"))
}
