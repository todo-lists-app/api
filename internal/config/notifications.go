package config

import (
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/caarlos0/env/v8"
)

type Notifications struct {
	VAPIDEmail   string `env:"VAPID_EMAIL" envDefault:""`
	VAPIDPrivate string `env:"VAPID_PRIVATE" envDefault:""`
	VAPIDPublic  string `env:"VAPID_PUBLIC" envDefault:""`
	TestUser     string `env:"NOTIFICATION_TEST_USER" envDefault:""`
}

func BuildNotifications(cfg *Config) error {
	notifications := &Notifications{}
	if err := env.Parse(notifications); err != nil {
		return logs.Errorf("unable to parse notifications: %w", err)
	}
	cfg.Notifications = *notifications

	return nil
}
