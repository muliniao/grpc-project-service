package grpc

import (
	"learning/grpc-project-service/pkg/util/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultPort = 3000
)

// Config Configuration for the GRPC Server Provider.
type Config struct {
	Port           int  // Port on which to start the GRPC service.
	LogInterceptor bool // Whether or not to enable logging interceptor.
	LogPayload     bool // Whether or not to enable logging of the payload. Should be disabled on production.
	EnableHealth   bool // Whether or not to register the health endpoint.

	// Json proto buffer marshaller config
	UseEnumAsInt        bool
	DisableEmitDefaults bool
}

// NewConfigFromEnv Initializes the configuration from environment variables or config file.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.AutomaticEnv()

	v.SetDefault("GRPC_PORT", defaultPort)
	v.SetDefault("GRPC_LOG_INTERCEPTOR", true)
	v.SetDefault("GRPC_LOG_PAYLOAD", false)
	v.SetDefault("GRPC_HEALTH_ENABLED", true)
	v.SetDefault("GRPC_USE_ENUM_AS_INT", false)
	v.SetDefault("GRPC_DISABLE_EMIT_DEFAULTS", false)

	// Load local config file
	config.LoadFromFile(v)

	port := v.GetInt("GRPC_PORT")
	logInterceptor := v.GetBool("GRPC_LOG_INTERCEPTOR")
	logPayload := v.GetBool("GRPC_LOG_PAYLOAD")
	enableHealth := v.GetBool("GRPC_HEALTH_ENABLED")
	useEnumAsInt := v.GetBool("GRPC_USE_ENUM_AS_INT")
	disableEmitDefaults := v.GetBool("GRPC_DISABLE_EMIT_DEFAULTS")

	logrus.WithFields(logrus.Fields{
		"port":         port,
		"logPayload":   logPayload,
		"enableHealth": enableHealth,
	}).Debug("Server Config Initialized")

	return &Config{
		Port:                port,
		LogInterceptor:      logInterceptor,
		LogPayload:          logPayload,
		EnableHealth:        enableHealth,
		UseEnumAsInt:        useEnumAsInt,
		DisableEmitDefaults: disableEmitDefaults,
	}
}
