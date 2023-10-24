package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

// StatusHandler is the handler for the status endpoint,
type StatusHandler struct {
	db *sqlx.DB
}

// NewStatusHandler initializes a new StatusHandler.
func NewStatusHandler(db *sqlx.DB) *StatusHandler {
	return &StatusHandler{
		db: db,
	}
}

// Status is the handler method for the status endpoint.
func (h *StatusHandler) Status(c *gin.Context) {
	db.Ping()
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
