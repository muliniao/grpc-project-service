package gateway

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"learning/grpc-project-service/pkg/util/config"
)

const (
	defaultPort = 8080
)

// Config configuration for the GRPC Gateway Provider.
type Config struct {
	Enabled    bool // Whether or not to enable the gateway.
	Port       int  // Port on which to start the HTTP service.
	LogPayload bool // Whether or not to enable logging of the payload. Should be disabled on production.
}

// NewConfigFromEnv initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("GRPC_GATEWAY_ENABLED", true)
	v.SetDefault("GRPC_GATEWAY_PORT", defaultPort)
	v.SetDefault("GRPC_GATEWAY_LOG_PAYLOAD", false)

	config.LoadFromFile(v)

	enabled := v.GetBool("GRPC_GATEWAY_ENABLED")
	port := v.GetInt("GRPC_GATEWAY_PORT")
	logPayload := v.GetBool("GRPC_GATEWAY_LOG_PAYLOAD")

	logrus.WithFields(logrus.Fields{
		"enabled":    enabled,
		"port":       port,
		"logPayload": logPayload,
	}).Debug("Gateway Config Initialized")

	return &Config{
		Enabled:    enabled,
		Port:       port,
		LogPayload: logPayload,
	}
}
