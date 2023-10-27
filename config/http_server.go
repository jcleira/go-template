package config

import "time"

type HTTPServer struct {
	ShutdownTimeout time.Duration `default:"10s"  envconfig:"HTTP_SERVER_SHUTDOWN_TIMEOUT"`
	Port            int           `default:"3000" envconfig:"HTTP_SERVER_PORT"`
}
