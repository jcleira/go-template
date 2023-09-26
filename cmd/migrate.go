package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	// This is needed to register the postgres driver with migrate.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

// MigrateCommand returns the cobra command to run migrations.
func MigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "runs migrations",
		Long:  "migrate will run migrations in the configured postgres instance",
		Run:   RunMigrations,
	}
}

// Settings define the configuration needed to run migrations.
type Settings struct {
	PostgresHost string `envconfig:"postgres_host"`
	PostgresPort string `envconfig:"postgres_port" default:"5432"`
	PostgresDB   string `envconfig:"postgres_db"`
	PostgresUser string `envconfig:"postgres_user"`
	PostgresPass string `envconfig:"postgres_password"`
}

// RunMigrations runs all the migrations.
func RunMigrations(_ *cobra.Command, _ []string) {
	var settings Settings
	if err := envconfig.Process("", &settings); err != nil {
		log.Fatalf("envconfig.Process. Err: %v", err)
	}

	migrateClient, err := migrate.New("file://data/db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			settings.PostgresUser, settings.PostgresPass,
			settings.PostgresHost, settings.PostgresPort,
			settings.PostgresDB),
	)
	if err != nil {
		log.Fatalf("migrate.New. Err: %v", err)
	}

	err = migrateClient.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrateClient.Up. Err: %v", err)
	}
}
