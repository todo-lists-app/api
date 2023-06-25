package api

import (
	"encoding/json"
	"errors"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"strings"
)

type Notification struct {
	config.Config
	context.Context
	UserID string
}

type UserSubscription struct {
	UserID      string
	LatestSub   string
	ChromeBased webpush.Subscription
	Mozilla     webpush.Subscription
	Apple       webpush.Subscription
	Edge        webpush.Subscription
	Unknown     webpush.Subscription
}

func NewNotificationService(ctx context.Context, cfg config.Config, id string) *Notification {
	return &Notification{
		Config:  cfg,
		Context: ctx,
		UserID:  id,
	}
}

func (n *Notification) StoreUser(subscription webpush.Subscription) error {
	sub := &UserSubscription{
		UserID: n.UserID,
	}

	subSet := false
	if strings.Contains(subscription.Endpoint, "googleapis.com") {
		sub.ChromeBased = subscription
		sub.LatestSub = "chromeBased"
		subSet = true
	}
	if strings.Contains(subscription.Endpoint, "mozilla.com") {
		sub.Mozilla = subscription
		sub.LatestSub = "mozilla"
		subSet = true
	}
	if strings.Contains(subscription.Endpoint, "apple.com") {
		sub.Apple = subscription
		sub.LatestSub = "apple"
		subSet = true
	}
	if strings.Contains(subscription.Endpoint, "edge.com") {
		sub.Edge = subscription
		sub.LatestSub = "edge"
		subSet = true
	}

	if !subSet {
		sub.Unknown = subscription
		sub.LatestSub = "unknown"
	}

	client, err := config.GetMongoClient(n.Context, n.Config)
	if err != nil {
		return logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(n.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	prev, err := n.GetSubscription(n.UserID)
	if err != nil {
		return logs.Errorf("error getting subscription: %v", err)
	}
	if prev != nil {
		if _, err := client.Database(n.Config.Mongo.Database).Collection(n.Config.Mongo.Collections.Notification).UpdateOne(n.Context,
			bson.D{{"$and", bson.A{
				bson.D{{"userID", n.UserID}},
				bson.D{{"development", n.Config.Local.Development}},
			}}},
			bson.D{{"$set", bson.M{
				"userID":       n.UserID,
				"latestSub":    sub.LatestSub,
				"developement": n.Config.Local.Development,
				sub.LatestSub:  subscription,
			}}}); err != nil {
			return logs.Errorf("error updating notification: %v", err)
		}

		return nil
	}

	if _, err := client.Database(n.Config.Mongo.Database).Collection(n.Config.Mongo.Collections.Notification).InsertOne(n.Context, bson.M{
		"userID":      n.UserID,
		"latestSub":   sub.LatestSub,
		"development": n.Config.Local.Development,
		sub.LatestSub: subscription,
	}); err != nil {
		return logs.Errorf("error inserting notification: %v", err)
	}

	return nil
}

func (n Notification) GetSubscription(userId string) (*UserSubscription, error) {
	client, err := config.GetMongoClient(n.Context, n.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(n.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	var sub UserSubscription
	if err := client.Database(n.Config.Mongo.Database).Collection(n.Config.Mongo.Collections.Notification).FindOne(n.Context, bson.M{
		"userID":       userId,
		"devvelopment": n.Config.Local.Development,
	}).Decode(&sub); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, logs.Errorf("error getting notification: %v", err)
		}
		return nil, nil
	}
	return &sub, nil
}

func (n Notification) SendTestNotification() error {
	client, err := config.GetMongoClient(n.Context, n.Config)
	if err != nil {
		return logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(n.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	var sub UserSubscription
	subscriber := n.Config.Notifications.VAPIDEmail
	if err := client.Database(n.Config.Mongo.Database).Collection(n.Config.Mongo.Collections.Notification).FindOne(n.Context, bson.M{
		"userID": "b3d1940e-d182-4fab-a574-37258e13d2d6",
	}).Decode(&sub); err != nil {
		return logs.Errorf("error getting notification: %v", err)
	}
	var useSub webpush.Subscription
	switch sub.LatestSub {
	case "chromeBased":
		useSub = sub.ChromeBased
		//subscriber = n.Config.Notifications.GoogleEmail
	case "mozilla":
		useSub = sub.Mozilla
	case "apple":
		useSub = sub.Apple
	case "edge":
		useSub = sub.Edge
	case "unknown":
		useSub = sub.Unknown
	}

	type MessageAction struct {
		Action string `json:"action"`
		Title  string `json:"title"`
	}
	type Message struct {
		Title   string          `json:"title"`
		Body    string          `json:"body"`
		Icon    string          `json:"icon"`
		Data    string          `json:"data"`
		Actions []MessageAction `json:"actions"`
	}

	m := Message{
		Title: "Test",
		Body:  "Test",
		Icon:  "https://beta.todo-list.app/logo512.png",
		Data:  "tm9kgx578a",
		Actions: []MessageAction{
			{
				Action: "open",
				Title:  "Go to Task",
			},
			{
				Action: "complete",
				Title:  "Task Complete",
			},
		},
	}
	message, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, err := webpush.SendNotification(message, &useSub, &webpush.Options{
		Subscriber:      subscriber,
		VAPIDPrivateKey: n.Notifications.VAPIDPrivate,
		VAPIDPublicKey:  n.Notifications.VAPIDPublic,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logs.Info(resp.StatusCode)
	logs.Infof("message: %+v", m)

	return nil
}
