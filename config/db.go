package config

import (
	"fmt"
	"net/url"
	"time"
)

// DB holds the configuration for the database.
type DB struct {
	Host string `required:"true"`
	Port int64  `required:"true"`
	User string `required:"true"`
	Pass string `required:"true"`
	Name string `required:"true"`

	// ConnectionTimeout is the maximum amount of time a connection may take to
	// establish.
	ConnectionTimeout time.Duration `default:"5s" envconfig:"CONNECTION_TIMEOUT"`
	// StatementTimeout is the maximum amount of time a statement may take to
	// execute.
	StatementTimeout time.Duration `default:"1m" envconfig:"STATEMENT_TIMEOUT"`
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle
	// before being closed.
	ConnMaxIdleTime time.Duration `default:"0" envconfig:"CONNECTION_MAX_IDLE_TIME"`
	// ConnMaxLifetime is the maximum amount of time a connection may exist
	// before being closed.
	ConnMaxLifetime time.Duration `default:"0" envconfig:"CONNECTION_MAX_LIFETIME"`

	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int `default:"20" envconfig:"MAX_OPEN_CONNS"`
	// MaxIdleConns is the maximum number of idle connections to the database.
	MaxIdleConns int `default:"20" envconfig:"MAX_IDLE_CONNS"`
	// SSLMode is the SSL mode to use when connecting to the database.
	SSLMode string `default:"require" envconfig:"DB_SSL_MODE"`
}

func (db DB) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&timezone=UTC&statement_timeout=%d",
		db.User, url.QueryEscape(db.Pass),
		db.Host, db.Port, db.Name,
		db.SSLMode, db.StatementTimeout.Milliseconds(),
	)
}
