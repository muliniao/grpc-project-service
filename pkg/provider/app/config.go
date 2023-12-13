package app

import (
	"encoding/json"
	"os"
	"path"
	"strings"

	"learning/grpc-project-service/pkg/util/config"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultBasePath = "/"
)

// Config configuration for the App Provider.
type Config struct {
	Name           string   // Application name.
	BasePath       string   // Base path.
	SpecialUrlPath []string // ignore basepath prefix, eg: ["/health", "/ping"]
}

// NewConfigFromEnv initializes the configuration from environment variables.
func NewConfigFromEnv() *Config {
	v := viper.New()
	v.AutomaticEnv()

	paths := strings.Split(os.Args[0], "/")
	v.SetDefault("APP_NAME", paths[len(paths)-1])
	v.SetDefault("APP_BASE_PATH", defaultBasePath)

	config.LoadFromFile(v)

	name := v.GetString("APP_NAME")
	basePath := path.Clean("/" + v.GetString("APP_BASE_PATH"))
	specialUrlPathBuf := v.GetString("SPECIAL_URL_PATH")
	specialUrlPath := make([]string, 0, 4)
	if len(specialUrlPathBuf) > 0 {
		err := json.Unmarshal([]byte(specialUrlPathBuf), &specialUrlPath)
		if err != nil {
			logrus.WithError(err).Debug("special_url_path error")
		}
	}

	logrus.WithFields(logrus.Fields{
		"name":             name,
		"base_path":        basePath,
		"special_url_path": specialUrlPathBuf,
	}).Debug("App Config initialized")

	return &Config{
		Name:           name,
		BasePath:       basePath,
		SpecialUrlPath: specialUrlPath,
	}
}
