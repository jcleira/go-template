package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// Healthz is the handler method for the status endpoint.
func (h *StatusHandler) Healthz(c *gin.Context) {
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
