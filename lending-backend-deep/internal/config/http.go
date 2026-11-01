package config

import (
	"fmt"
	"time"
)

// EnvHTTP holds environment variables that configure the HTTP server.
type EnvHTTP struct {
	HTTPServerHost            string `env:"HTTP_SERVER_HOST"             envDefault:"0.0.0.0"`
	HTTPServerPort            int    `env:"HTTP_SERVER_PORT"             envDefault:"8080"`
	HTTPServerReadTimeout     int    `env:"HTTP_SERVER_READ_TIMEOUT"     envDefault:"30"`
	HTTPServerWriteTimeout    int    `env:"HTTP_SERVER_WRITE_TIMEOUT"    envDefault:"30"`
	HTTPServerShutdownTimeout int    `env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" envDefault:"5"`
}

// ServerAddr returns the "host:port" string the HTTP server binds to.
func (e *EnvHTTP) ServerAddr() string {
	return fmt.Sprintf("%s:%d", e.HTTPServerHost, e.HTTPServerPort)
}

// ReadTimeout is the maximum duration for reading the entire request.
func (e *EnvHTTP) ReadTimeout() time.Duration {
	return time.Duration(e.HTTPServerReadTimeout) * time.Second
}

// WriteTimeout is the maximum duration before timing out writes of the response.
func (e *EnvHTTP) WriteTimeout() time.Duration {
	return time.Duration(e.HTTPServerWriteTimeout) * time.Second
}

// ShutdownTimeout is the maximum duration allowed for a graceful shutdown.
func (e *EnvHTTP) ShutdownTimeout() time.Duration {
	return time.Duration(e.HTTPServerShutdownTimeout) * time.Second
}
