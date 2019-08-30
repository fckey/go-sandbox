package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/fckey/go-sandbox/pubsubop"
	"github.com/fckey/go-sandbox/slackop"
	"github.com/fckey/go-sandbox/twiop"
)

var (
	mgr    *twiop.Manager
	strmgr *twiop.StreamManager
	psMgr  *pubsubop.Manager
	since  *time.Time
	slack  *slackop.Manager
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "World"
	}

	fmt.Fprintf(w, "Hello %s: %s!\n", target, "")
}

func polling(w http.ResponseWriter, r *http.Request) {
	log.Print("Polling Tweets.")
	testSearch()
	fmt.Fprintf(w, "Polling done!\n")
}

func stop(w http.ResponseWriter, r *http.Request) {
	log.Print("Stop Stream.")
	strmgr.StopStream()
	fmt.Fprintf(w, "Stream Stopped!\n")
}

func createSub(w http.ResponseWriter, r *http.Request) {
	err := psMgr.CreateSub()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, fmt.Sprintf("Created ", psMgr.SubName))
}

func getPubSub(w http.ResponseWriter, r *http.Request) {
	messages, err := psMgr.PullMessages()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	log.Println(strings.Join(messages[:], ","))
	fmt.Fprintf(w, "Logged extracted message")
}

func ping(w http.ResponseWriter, r *http.Request) {
	if since == nil {
		log.Println("Setting since")
		tmp := time.Now()
		tmp = tmp.Add(-time.Hour)
		since = &tmp
	}
	f := "2006-01-02_15:04:05_MST"
	fmt.Fprintf(w, "Ack "+since.Format(f))
}

func testSearch() {
	if since == nil {
		tmp := time.Now()
		tmp = tmp.Add(-time.Hour)
		since = &tmp
	}
	tweets := mgr.FilterSearch(mgr.SearchKeywordSince("Cloud Google", *since), func(t twitter.Tweet) bool {
		return len(t.Text) > 0 && t.InReplyToStatusID == 0
	})
	for _, t := range tweets {
		log.Println(fmt.Printf("Working for %v\n", t))
		bytes, err := json.Marshal(twiop.Simple(t))
		if err != nil {
			log.Fatal(err)
		}
		psMgr.Publish(bytes)
	}
	now := time.Now()
	since = &now
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", handler)
	http.HandleFunc("/polling", polling)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/createsub", createSub)
	http.HandleFunc("/getpubsub", getPubSub)
	http.HandleFunc("/ping", ping)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mgr = &twiop.Manager{}
	strmgr = &twiop.StreamManager{}
	mgr.InitClient()

	projectID := os.Getenv("GCP_PROJECT_ID")
	topicName := os.Getenv("PUBSUB_TOPIC")
	twiSub := os.Getenv("PUBSUB_SUBSCRIPTION")

	psMgr = &pubsubop.Manager{
		ProjectID: projectID,
		TopicName: topicName,
		SubName:   twiSub,
	}
	psMgr.InitClient()

	url := os.Getenv("SLACK_NOTIFY_URL")
	slack = &slackop.Manager{
		URL: url,
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
