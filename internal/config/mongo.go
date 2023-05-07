package config

import (
	env "github.com/caarlos0/env/v8"
	vh "github.com/keloran/vault-helper"
)

type DB struct {
	Database       string
	KeysCollection string
}

type Mongo struct {
	Host           string `env:"MONGO_HOST" envDefault:"localhost"`
	Username       string `env:"MONGO_USER" envDefault:""`
	Password       string `env:"MONGO_PASS" envDefault:""`
	Database       string `env:"MONGO_DB" envDefault:""`
	KeyCollection  string `env:"MONGO_KEY_COLLECTION" envDefault:""`
	ListCollection string `env:"MONGO_LIST_COLLECTION" envDefault:""`
	MongoPath      string `env:"MONGO_VAULT_PATH" envDefault:""`
}

func BuildMongo(c *Config) error {
	mongo := &Mongo{}

	if err := env.Parse(mongo); err != nil {
		return err
	}

	v := vh.NewVault(c.Vault.Address, c.Vault.Token)
	if err := v.GetSecrets(c.MongoPath); err != nil {
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

	mongo.Password = password
	mongo.Username = username

	c.Mongo = *mongo

	return nil
}
