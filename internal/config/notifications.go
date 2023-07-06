package config

import "github.com/caarlos0/env/v8"

type Notifications struct {
	VAPIDEmail   string `env:"VAPID_EMAIL" envDefault:"" json:"vapid_email,omitempty"`
	VAPIDPrivate string `env:"VAPID_PRIVATE" envDefault:"" json:"vapid_private,omitempty"`
	VAPIDPublic  string `env:"VAPID_PUBLIC" envDefault:"" json:"vapid_public,omitempty"`
	TestUser     string `env:"NOTIFICATION_TEST_USER" envDefault:"" json:"test_user,omitempty"`
}

func BuildNotifications(cfg *Config) error {
	notifications := &Notifications{}
	if err := env.Parse(notifications); err != nil {
		return err
	}
	cfg.Notifications = *notifications

	return nil
}
