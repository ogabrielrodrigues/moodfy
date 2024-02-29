package style

import (
	"database/sql"
	"net/http"
)

type handler struct {
	db *sql.DB
}

func Handler(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) CreateStyle(w http.ResponseWriter, r *http.Request) {}

func (h *handler) ListStyle(w http.ResponseWriter, r *http.Request) {}
