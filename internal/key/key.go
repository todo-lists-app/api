package key

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/scrypt"
)

type ExportKey struct {
	DecodeKey string `json:"decodekey,omitempty" bson:"decodekey,omitempty"`
	UserId    string `json:"userid,omitempty" bson:"userid,omitempty"`
}

type Key struct {
	config.Config
	context.Context
}

func NewKeyService(ctx context.Context, cfg config.Config) *Key {
	return &Key{
		Config:  cfg,
		Context: ctx,
	}
}

func (k *Key) GetKey(userid string) (*ExportKey, error) {
	var storedKey ExportKey

	client, err := mongo.Connect(k.Context, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", k.Config.Mongo.Username, k.Config.Mongo.Password, k.Config.Mongo.Host)))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(k.Context)

	if err := client.Database(k.Config.Mongo.Database).Collection(k.Config.Mongo.KeyCollection).FindOne(k.Context, &bson.M{
		"userid": userid,
	}).Decode(&storedKey); err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	return &storedKey, nil
}

func (k *Key) CreateKey(id, salt string) (*ExportKey, error) {
	dk, err := scrypt.Key([]byte(id), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(k.Context, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", k.Config.Mongo.Username, k.Config.Mongo.Password, k.Config.Mongo.Host)))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(k.Context)

	if _, err := client.Database(k.Config.Mongo.Database).Collection(k.Config.Mongo.KeyCollection).InsertOne(k.Context, &bson.M{
		"userid":    id,
		"decodekey": base64.StdEncoding.EncodeToString(dk),
	}); err != nil {
		return nil, err
	}

	return &ExportKey{
		DecodeKey: base64.StdEncoding.EncodeToString(dk),
		UserId:    id,
	}, nil
}
