package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"

	"github.com/jcleira/go-template/config"
)

var Web = &cobra.Command{
	Use:   "web",
	Short: "web starts running the http server",
	RunE:  RunWeb,
}

func RunWeb(_ *cobra.Command, _ []string) error {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	config, err := config.Get()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// Open the database connection.
	_, err := sqlx.Open("postgres", config.DB.URL())
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	r := gin.Default()

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server.ListenAndServe", slog.String("error", err.Error()))
		}
	}()

	ctx, cancel := context.WithTimeout(
		context.Background(), config.HTTPServer.ShutdownTimeout)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("Received signal", slog.String("signal", sig.String()))

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server.Shutdown", slog.String("error", err.Error()))
		return err
	}

	return nil
}
