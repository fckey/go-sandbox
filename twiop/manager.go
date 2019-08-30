package twiop

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type SimpleTweet struct {
	CreatedAt            string `json:"created_at"`
	FavoriteCount        int    `json:"favorite_count"`
	Favorited            bool   `json:"favorited"`
	InReplyToScreenName  string `json:"in_reply_to_screen_name"`
	InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
	InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
	InReplyToUserID      int64  `json:"in_reply_to_user_id"`
	InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
	Lang                 string `json:"lang"`
	PossiblySensitive    bool   `json:"possibly_sensitive"`
	QuoteCount           int    `json:"quote_count"`
	ReplyCount           int    `json:"reply_count"`
	RetweetCount         int    `json:"retweet_count"`
	Retweeted            bool   `json:"retweeted"`
	Text                 string `json:"text"`
	FullText             string `json:"full_text"`
	Truncated            bool   `json:"truncated"`
	//User                 *User                  `json:"user"`
	ID         int64  `json:"id"`
	IDStr      string `json:"id_str"`
	ScreenName string `json:"screen_name"`
	Timezone   string `json:"time_zone"`
	// end user
	QuotedStatusID    int64  `json:"quoted_status_id"`
	QuotedStatusIDStr string `json:"quoted_status_id_str"`
}

// Manager is wrapper of twitter interface
type Manager struct {
	client *twitter.Client
}

// StreamManager is wrapper of twitter streaming interface
type StreamManager struct {
	client *twitter.Client
	stream *twitter.Stream
}

func createHTTPClient() *http.Client {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	return config.Client(oauth1.NoContext, token)
}

// InitClient initialize twitter client
func (mgr *Manager) InitClient() {
	httpClient := createHTTPClient()
	mgr.client = twitter.NewClient(httpClient)
}

// InitClient initialize twitter client for streaming
func (mgr *StreamManager) InitClient() {
	httpClient := createHTTPClient()
	mgr.client = twitter.NewClient(httpClient)
}

// TrackKeywords get tweets realtime of a keyword
func (mgr *StreamManager) TrackKeywords(process func(tweet *twitter.Tweet), keyword string) {
	demux := twitter.NewSwitchDemux()

	if process == nil {
		demux.Tweet = func(tweet *twitter.Tweet) {
			fmt.Println(tweet.Text)
		}
	} else {
		demux.Tweet = process
	}

	fmt.Println("Starting Stream...")

	filterParams := &twitter.StreamFilterParams{
		Track:         []string{keyword},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := mgr.client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}
	mgr.stream = stream
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
}

// StopStream stops active stream
func (mgr *StreamManager) StopStream() {
	if mgr.stream != nil {
		fmt.Println("Stopping Stream...")
		mgr.stream.Stop()
	} else {
		log.Println("No actice Stream.")
	}
}

// SearchKeyword looks up keyword in twitter
func (mgr *Manager) SearchKeywordSince(keyword string, since time.Time) (search *twitter.Search) {
	f := "2006-01-02_15:04:05_MST"

	search, response, err := mgr.client.Search.Tweets(&twitter.SearchTweetParams{
		Query: keyword,
		Since: since.Format(f),
	})
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Response code is %v\n", response.StatusCode)
	}
	return
}

// SearchKeyword looks up keyword in twitter
func (mgr *Manager) SearchKeyword(keyword string) (search *twitter.Search) {
	search, response, err := mgr.client.Search.Tweets(&twitter.SearchTweetParams{
		Query: keyword,
		Since: "2019-08-12",
	})
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Response code is %v\n", response.StatusCode)
	}
	return
}

// Mention create mention to tweet
func (mgr *Manager) Mention(t twitter.Tweet, msg string) {
	// id := t.User.ID
	screenName := t.User.ScreenName
	reply := fmt.Sprintf("@%v %s", screenName, msg)

	p := twitter.StatusUpdateParams{
		InReplyToStatusID: t.ID,
	}
	log.Println("Send: ", reply)
	_, response, err := mgr.client.Statuses.Update(reply, &p)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("Response code is %v \n", response.StatusCode)
	}
}

// FilterSearch collects tweets matching with filter
// f should validate if the tweet is valid to reply such as having specific keyword
func (mgr *Manager) FilterSearch(s *twitter.Search, f func(t twitter.Tweet) bool) []twitter.Tweet {
	tweets := []twitter.Tweet{}
	for _, t := range s.Statuses {
		if f(t) {
			tweets = append(tweets, t)
			log.Println(t.Text)
		}
	}
	return tweets
}

func Simple(t twitter.Tweet) SimpleTweet {
	return SimpleTweet{
		CreatedAt:            t.CreatedAt,
		FavoriteCount:        t.FavoriteCount,
		Favorited:            t.Favorited,
		InReplyToScreenName:  t.InReplyToScreenName,
		InReplyToStatusID:    t.InReplyToStatusID,
		InReplyToStatusIDStr: t.InReplyToStatusIDStr,
		InReplyToUserID:      t.InReplyToUserID,
		InReplyToUserIDStr:   t.InReplyToUserIDStr,
		Lang:                 t.Lang,
		PossiblySensitive:    t.PossiblySensitive,
		QuoteCount:           t.QuoteCount,
		ReplyCount:           t.ReplyCount,
		RetweetCount:         t.RetweetCount,
		Retweeted:            t.Retweeted,
		Text:                 t.Text,
		FullText:             t.FullText,
		Truncated:            t.Truncated,
		//User                 *User                  `json:"user"`
		ID:         t.User.ID,
		IDStr:      t.User.IDStr,
		ScreenName: t.User.ScreenName,
		Timezone:   t.User.Timezone,
		// end user
		QuotedStatusID:    t.QuotedStatusID,
		QuotedStatusIDStr: t.QuotedStatusIDStr,
	}
}
