package slackop

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// payload to send to Slack API
type payload struct {
	text string `json:"text"`
	name string `json:name`
}

// Target contains Slack API endpoint and message to be sent
type Target struct {
	URL  string
	text string
}

// Manager keeps context to interact with slack
type Manager struct {
	URL string
}

// Notify send text to URL in Target
func Notify(tg Target) error {
	msg := payload{
		text: tg.text,
	}
	m, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("Failed to marchall msg: ", err)
		return err
	}
	req, err := http.NewRequest(
		"POST",
		tg.URL,
		bytes.NewBuffer(m),
	)
	if err != nil {
		log.Fatal("Failed to create request: ", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	resp, err := c.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Failed to POST msg: ", err)
		return err
	}
	return nil
}

func (mgr *Manager) Notify(text string) error {
	return Notify(Target{URL: mgr.URL, text: text})
}
