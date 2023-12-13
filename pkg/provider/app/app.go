// Package app application Provider.
// Maintains basic info for the application (like name and version).
package app

import (
	"path"

	"learning/grpc-project-service/pkg/provider"
	"learning/grpc-project-service/pkg/version"

	"github.com/sirupsen/logrus"
)

// App application Provider.
// Maintains basic info for the application (like name and version).
type App struct {
	provider.AbstractProvider

	Config *Config
}

// New creates an App Provider.
func New(config *Config) *App {
	if config == nil {
		config = NewConfigFromEnv()
	}

	return &App{
		Config: config,
	}
}

// Init app provider doesn't need initialization, since version is set during compilation and name via environment variables.
func (p *App) Init() error {
	logrus.WithFields(logrus.Fields{
		"name":    p.Name(),
		"version": p.Version().String(),
	}).Info("App Provider initialized")
	return nil
}

// Name returns the application name.
func (p *App) Name() string {
	return p.Config.Name
}

// SpecialUrlPath returns the SpecialUrlPath, eg: ["/health", "/ping"].
func (p *App) SpecialUrlPath() []string {
	return p.Config.SpecialUrlPath
}

// Version returns the Application version.
func (p *App) Version() version.Version {
	return version.CurrentVersion()
}

// ParsePath appends the given elements to the base path and returns a cleaned URL path.
// The resulting path will always end with a "/".
func (p *App) ParsePath(elem ...string) string {
	res := p.ParseEndpoint(elem...)
	if res != "/" {
		res += "/"
	}
	return res
}

// ParseEndpoint appends the given elements to the base path and returns a cleaned URL path.
// The resulting path will not end with a "/", unless that's the only character it contains (root path).
func (p *App) ParseEndpoint(elem ...string) string {
	elem = append([]string{p.Config.BasePath}, elem...)
	return path.Join(elem...)
}
