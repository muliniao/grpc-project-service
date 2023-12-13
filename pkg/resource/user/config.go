package user

import "github.com/spf13/viper"

type Config struct {
	UserAddr    string
	MockEnabled bool
}

const (
	defaultUserMockEnabled = true
	defaultBillingAddr     = "127.0.0.1:3000"
)

func NewConfigFromEnv() *Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("USER")

	v.SetDefault("ADDR", defaultBillingAddr)
	v.SetDefault("MOCK_ENABLED", defaultUserMockEnabled)
	return &Config{
		UserAddr:    v.GetString("ADDR"),
		MockEnabled: v.GetBool("MOCK_ENABLED"),
	}
}
