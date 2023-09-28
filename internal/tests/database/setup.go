// Package database provides database setup functions for testing.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"

	"github.com/jcleira/go-template/config"
)

const (
	databaseName = "alerts-test"
	databaseUser = "alerts-test"
	databasePass = "alerts-test"

	imageName    = "postgres"
	imageVersion = "latest"

	dockerAuthEndpoint           = "https://index.docker.io/v1/"
	dockerResourceExpirationTime = 300
)

func CreateTestingDatabase(t *testing.T) (db *sql.DB, teardown func()) {
	t.Helper()

	pool, err := dockertest.NewPool("")
	require.NoError(t, err)

	resource, err := createDockerResource(pool)
	require.NoError(t, err)

	db, err = waitTillConnection(pool, resource)
	require.NoError(t, err)

	require.NoError(t, runMigrations(db, "../../../../data/db/migrations"))

	return db, func() {
		require.NoError(t, pool.Purge(resource))
	}
}

func createDockerResource(pool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: imageName,
			Tag:        imageVersion,
			Env: []string{
				"POSTGRES_DB=" + databaseName,
				"POSTGRES_USER=" + databaseUser,
				"POSTGRES_PASSWORD=" + databasePass,
				"PGTZ=UTC",
			},
			Cmd: []string{
				"postgres",
				"-c",
				"log_statement=all",
			},
		},

		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return nil, fmt.Errorf("pool.RunWithOptions: %w", err)
	}

	if err = resource.Expire(dockerResourceExpirationTime); err != nil {
		return nil, fmt.Errorf("resource.Expire: %w", err)
	}

	return resource, nil
}

func waitTillConnection(
	pool *dockertest.Pool, resource *dockertest.Resource) (*sql.DB, error) {
	var db *sql.DB

	// The host is different between CI and locally
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	portString := resource.GetPort("5432/tcp")
	port, err := strconv.ParseInt(portString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	configDB := config.DB{
		Host:    host,
		Port:    port,
		User:    databaseUser,
		Pass:    databasePass,
		Name:    databaseName,
		SSLMode: "disable",
	}

	fmt.Println(configDB.URL())

	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", configDB.URL())
		if err != nil {
			return fmt.Errorf("open database connection: %w", err)
		}

		if err := db.Ping(); err != nil {
			return fmt.Errorf("ping database: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("pool.Retry: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath), "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("m.Up: %w", err)
	}

	return nil
}
