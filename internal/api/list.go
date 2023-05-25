package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// List is the list service
type List struct {
	config.Config
	context.Context
	UserID string
}

// NewListService creates a new list service
func NewListService(ctx context.Context, cfg config.Config, id string) *List {
	return &List{
		Config:  cfg,
		Context: ctx,
		UserID:  id,
	}
}

// StoredList is the stored list
type StoredList struct {
	UserID string `bson:"userid" json:"userid"`
	Data   string `bson:"data" json:"data"`
	IV     string `bson:"iv" json:"iv"`
}

// GetList gets a list for the user
func (l *List) GetList() (*StoredList, error) {
	var storeList StoredList

	client, err := mongo.Connect(l.Context, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", l.Config.Mongo.Username, l.Config.Mongo.Password, l.Config.Mongo.Host)))
	if err != nil {
		return &storeList, err
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Debugf("error disconnecting from mongo: %v", err)
		}
	}()

	if err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.ListCollection).FindOne(l.Context, &bson.M{
		"userid": l.UserID,
	}).Decode(&storeList); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &storeList, err
		}
		return &storeList, nil
	}

	return &storeList, nil
}

// UpdateList updates a list for the user
func (l *List) UpdateList(list *StoredList) (*StoredList, error) {
	// if _, err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.ListCollection).UpdateOne(l.Context, bson.M{
	//	"userid": l.UserID,
	// }, bson.D{
	//	{
	//		"$set", bson.D{
	//			{"userid", l.UserID},
	//			{"data", list.Data},
	//			{"iv", list.IV},
	//		},
	//	},
	// }); err != nil {
	//	return nil, err
	// }

	return list, nil
}

// DeleteList deletes a list for the user
func (l *List) DeleteList(id string) (*StoredList, error) {
	return &StoredList{
		UserID: id,
	}, nil
}

// CreateList creates a new list for the user
func (l *List) CreateList(list *StoredList) (*StoredList, error) {
	client, err := mongo.Connect(l.Context, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", l.Config.Mongo.Username, l.Config.Mongo.Password, l.Config.Mongo.Host)))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Debugf("error disconnecting from mongo: %v", err)
		}
	}()

	if _, err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.ListCollection).InsertOne(l.Context, bson.M{
		"userid": l.UserID,
		"data":   list.Data,
		"iv":     list.IV,
	}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return list, nil
		}
		return nil, err
	}

	return list, nil
}
