// Package config provides a struct to store the applications config
package config

import (
	"github.com/metal-toolbox/rivets/ginjwt"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/otelx"
)

// AppConfig stores all the config values for our application
var AppConfig struct {
	CRDB    crdbx.Config
	Logging loggingx.Config
	Tracing otelx.Config
	// APIServerJWTAuth sets the JWT verification configuration for the conditionorc API service.
	APIServerJWTAuth []ginjwt.AuthConfig `mapstructure:"ginjwt_auth"`
}
