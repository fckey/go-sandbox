package pubsubop

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
)

type Manager struct {
	ProjectID string
	SubName   string
	TopicName string
	topic     *pubsub.Topic
	client    *pubsub.Client
}

func (mgr *Manager) InitClient() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, mgr.ProjectID)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	mgr.client = client
}

func (mgr *Manager) PullMessages() ([]string, error) {
	ctx := context.Background()
	messages := []string{}

	var mu sync.Mutex
	received := 0
	sub := mgr.client.Subscription(mgr.SubName)
	cctx, cancel := context.WithCancel(ctx)

	err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		msg.Ack()
		m := string(msg.Data)
		messages = append(messages, m)
		mu.Lock()
		defer mu.Unlock()
		received++
		if received >= 5 {
			cancel()
		}
	})

	return messages, err
}

// PublishStr sends data given as string to pubsub server
func (mgr *Manager) PublishStr(msg string) error {
	return mgr.Publish([]byte(msg))
}

// Publish sends data given as array of byte to pubsub server
func (mgr *Manager) Publish(msg []byte) error {
	ctx := context.Background()
	t := mgr.getTopic()
	log.Printf(string(msg))
	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})
	// Block until the result is returned and id server-generated
	// ID is for the published message
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}

// CreateSub creates subscription
func (mgr *Manager) CreateSub() error {
	t := mgr.getTopic()
	err := mgr.createSub(mgr.SubName, t)
	if err != nil {
		return err
	}
	return nil
}

func (mgr *Manager) createSub(subName string, topic *pubsub.Topic) error {
	ctx := context.Background()
	sub, err := mgr.client.CreateSubscription(ctx, subName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 30 * time.Second,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Created subscription: %v\n", sub)
	// [END pubsub_create_pull_subscription]
	return nil
}

func (mgr *Manager) getTopic() *pubsub.Topic {
	if mgr.topic != nil {
		return mgr.topic
	}
	t := mgr.createTopicIfNotExists(mgr.TopicName)
	mgr.topic = t
	return t
}

func (mgr *Manager) createTopicIfNotExists(topicName string) *pubsub.Topic {
	ctx := context.Background()

	t := mgr.client.Topic(topicName)

	ok, err := t.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		return t
	}

	t, err = mgr.client.CreateTopic(ctx, topicName)
	if err != nil {
		log.Fatalf("Failed to create the TopicName: %v", err)
	}
	return t
}
