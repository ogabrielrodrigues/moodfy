package artist

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
)

type handler struct {
	db *sql.DB
}

func Handler(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) CreateArtist(w http.ResponseWriter, r *http.Request) {
	body := DTO{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "corpo da requisição malformado ou inválido",
		})
		return
	}

	if len(body.Name) < 3 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "o nome do artista tem que ter ao menos 3 caracteres",
		})
		return
	}

	ctx := context.Background()
	conn, err := h.db.Conn(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao se conectar ao banco de dados",
		})
		return
	}

	artist := New(body.Name)

	row := conn.QueryRowContext(ctx, `
		INSERT INTO "artist" (name)
		VALUES ($1) RETURNING id`,
		artist.Name,
	)

	if err := row.Scan(&artist.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao inserir o registro no banco de dados",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

func (h *handler) ListArtist(w http.ResponseWriter, r *http.Request) {}
