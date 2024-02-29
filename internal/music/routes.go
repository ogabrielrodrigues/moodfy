package music

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

func (h *handler) CreateMusic(w http.ResponseWriter, r *http.Request) {}

// List or filter musics
func (h *handler) ListMusic(w http.ResponseWriter, r *http.Request) {}
