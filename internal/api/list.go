package api

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"github.com/todo-lists-app/todo-lists-api/internal/key"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoList struct {
	UserID    string `json:"userid,omitempty"`
	Created   string `json:"created,omitempty"`
	Updated   string `json:"updated,omitempty"`
	DecodeKey string `json:"decodekey,omitempty"`
	Data      string `json:"data,omitempty"`
}

type List struct {
	config.Config
	context.Context
	UserID string
	Key    string
}

func NewListService(ctx context.Context, cfg config.Config, id, key string) *List {
	return &List{
		Config:  cfg,
		Context: ctx,
		UserID:  id,
		Key:     key,
	}
}

type StoredList struct {
	UserID string `bson:"userid"`
	Data   string `bson:"data"`
}

func (l *List) GetList() (*TodoList, error) {
	var dc key.ExportKey

	decodeKey, err := key.NewKeyService(l.Context, l.Config).GetKey(l.UserID)
	if err != nil {
		return nil, err
	}

	if decodeKey != nil && decodeKey.DecodeKey != "" {
		dc = *decodeKey
	} else {
		decodeKey, err := key.NewKeyService(l.Context, l.Config).CreateKey(l.UserID, l.Key)
		if err != nil {
			return nil, err
		}
		if decodeKey != nil && decodeKey.DecodeKey != "" {
			dc = *decodeKey
		} else {
			return nil, logs.Local().Errorf("unable to create key")
		}
	}

	client, err := mongo.Connect(l.Context, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", l.Config.Mongo.Username, l.Config.Mongo.Password, l.Config.Mongo.Host)))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(l.Context)

	var storeList StoredList
	if err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.ListCollection).FindOne(l.Context, &bson.M{
		"userid": l.UserID,
	}).Decode(&storeList); err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	return &TodoList{
		DecodeKey: dc.DecodeKey,
		Data:      storeList.Data,
	}, nil
}

func (l *List) UpdateList(list *TodoList) (*TodoList, error) {
	return &TodoList{}, nil
}

func (l *List) DeleteList(id string) (*TodoList, error) {
	return &TodoList{}, nil
}

func (l *List) CreateList(list *TodoList) (*TodoList, error) {
	return &TodoList{}, nil
}
