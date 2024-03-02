package style

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type handler struct {
	repo repo
}

func Handler(db *sql.DB) *handler {
	repo := *Repo(db)
	return &handler{repo}
}

func (h *handler) CreateStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body := DTO{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "corpo da requisição malformado ou inválido",
		})
		return
	}

	style, code, err := h.repo.Create(&body)
	if err != nil {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(code),
			"message": err.Error(),
		})
		return
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(style)
}

func (h *handler) ListStyle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	styles, code, err := h.repo.List()
	if err != nil {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(code),
			"message": err.Error(),
		})
		return
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(styles)
}
