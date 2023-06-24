package api

import (
	"context"
	"errors"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	client, err := config.GetMongoClient(l.Context, l.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	storedList := StoredList{}
	if err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.Collections.List).FindOne(l.Context, &bson.M{
		"userid": l.UserID,
	}).Decode(&storedList); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return &storedList, logs.Errorf("error finding list: %v", err)
		}
		return &storedList, nil
	}

	return &storedList, nil
}

// UpdateList updates a list for the user
func (l *List) UpdateList(list *StoredList) (*StoredList, error) {
	client, err := config.GetMongoClient(l.Context, l.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	if _, err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.Collections.List).UpdateOne(l.Context,
		bson.D{{"userid", l.UserID}},
		bson.D{{"$set", bson.M{
			"data": list.Data,
			"iv":   list.IV,
		}}}); err != nil {
		return nil, logs.Errorf("error updating list: %v", err)
	}

	return list, nil
}

// DeleteList deletes a list for the user
func (l *List) DeleteList(id string) (*StoredList, error) {
	client, err := config.GetMongoClient(l.Context, l.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	return &StoredList{
		UserID: id,
	}, nil
}

// CreateList creates a new list for the user
func (l *List) CreateList(list *StoredList) (*StoredList, error) {
	client, err := config.GetMongoClient(l.Context, l.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	if _, err := client.Database(l.Config.Mongo.Database).Collection(l.Config.Mongo.Collections.List).InsertOne(l.Context, bson.M{
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
