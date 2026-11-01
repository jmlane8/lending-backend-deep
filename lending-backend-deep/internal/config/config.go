package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
)

// LoadFromEnv parses all environment variables into an Env struct.
// Panics on parse failure so misconfiguration is caught at startup.
func LoadFromEnv() *Env {
	var cfg Env
	if err := env.Parse(&cfg); err != nil {
		panic(fmt.Sprintf("config: failed to parse env vars: %v", err))
	}
	return &cfg
}

// Env is the root configuration struct.  Sub-structs are composed for
// readability.  All fields are populated from environment variables.
type Env struct {
	EnvHTTP
	EnvPostgreSQL
	EnvWorker

	AppName string `env:"APP_NAME" envDefault:"lending-backend-deep"`
	AppEnv  string `env:"APP_ENV" envDefault:"local"` // local|dev|uat|sit|sandbox|prod

	LoggerLevel  string `env:"LOGGER_LEVEL"  envDefault:"info"`
	LoggerOutput string `env:"LOGGER_OUTPUT" envDefault:"stdout"` // console|stdout

	// OpenAPIUIEnable enables serving the bundled OpenAPI spec and Redoc UI.
	OpenAPIUIEnable bool `env:"OPENAPI_UI_ENABLE" envDefault:"false"`

	// OpenAPIUIPath is the URL path at which the Redoc UI is served.
	OpenAPIUIPath string `env:"OPENAPI_UI_PATH"   envDefault:"/redoc"`
}

// Logger builds a zerolog.Logger configured by LoggerLevel and LoggerOutput.
func (e *Env) Logger() zerolog.Logger {
	level, err := zerolog.ParseLevel(e.LoggerLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	if e.LoggerOutput == "console" {
		return zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stdout
		})).Level(level).With().Timestamp().Logger()
	}

	return zerolog.New(os.Stdout).
		Level(level).
		With().
		Timestamp().
		Str("app", e.AppName).
		Str("env", e.AppEnv).
		Logger()
}
