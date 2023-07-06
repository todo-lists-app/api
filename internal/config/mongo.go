package config

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	env "github.com/caarlos0/env/v8"
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
	mongo := &Mongo{}

	if err := env.Parse(mongo); err != nil {
		return err
	}

	v := vh.NewVault(c.Vault.Address, c.Vault.Token)
	if err := v.GetSecrets(mongo.Vault.Path); err != nil {
		return err
	}

	username, err := v.GetSecret("username")
	if err != nil {
		return err
	}

	password, err := v.GetSecret("password")
	if err != nil {
		return err
	}

	mongo.Vault.ExpireTime = time.Now().Add(time.Duration(v.LeaseDuration) * time.Second)
	mongo.Password = password
	mongo.Username = username

	c.Mongo = *mongo

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
