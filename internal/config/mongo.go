package config

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
	vh "github.com/keloran/vault-helper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// Mongo is the Mongo config
type Mongo struct {
	Host     string `env:"MONGO_HOST" envDefault:"localhost"`
	Username string `env:"MONGO_USER" envDefault:""`
	Password string `env:"MONGO_PASS" envDefault:""`
	Database string `env:"MONGO_DB" envDefault:""`

	Collections struct {
		Account      string `env:"MONGO_ACCOUNT_COLLECTION" envDefault:""`
		List         string `env:"MONGO_LIST_COLLECTION" envDefault:""`
		Notification string `env:"MONGO_NOTIFICATION_COLLECTION" envDefault:""`
	}
	Vault struct {
		Path       string `env:"MONGO_VAULT_PATH" envDefault:""`
		ExpireTime time.Time
	}
}

// BuildMongo builds the Mongo config
func BuildMongo(c *Config) error {
	mungo := &Mongo{}

	if err := env.Parse(mungo); err != nil {
		return logs.Errorf("error parsing mongo: %v", err)
	}

	v := vh.NewVault(c.Vault.Address, c.Vault.Token)
	if err := v.GetSecrets(mungo.Vault.Path); err != nil {
		return logs.Errorf("error getting secrets: %v", err)
	}

	username, err := v.GetSecret("username")
	if err != nil {
		return logs.Error("error getting username")
	}

	password, err := v.GetSecret("password")
	if err != nil {
		return logs.Errorf("error getting password: %v", err)
	}

	mungo.Vault.ExpireTime = time.Now().Add(time.Duration(v.LeaseDuration) * time.Second)
	mungo.Password = password
	mungo.Username = username

	c.Mongo = *mungo

	return nil
}

func GetMongoClient(ctx context.Context, cfg Config) (*mongo.Client, error) {
	if time.Now().Unix() > cfg.Mongo.Vault.ExpireTime.Unix() {
		if err := BuildMongo(&cfg); err != nil {
			return nil, logs.Errorf("error re-building mongo: %v", err)
		}
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host)))
	if err != nil {
		return nil, logs.Errorf("error connecting to mongo: %v", err)
	}

	return client, nil
}
