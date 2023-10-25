package config

import "time"

type HTTPServer struct {
	ShutdownTimeout time.Duration `envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT" default:"10s"`
	Port            int           `envconfig:"HTTP_SERVER_PORT" default:"3000"`
}
