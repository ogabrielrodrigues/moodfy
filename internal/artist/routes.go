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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusBadRequest),
			"message": "corpo da requisição malformado ou inválido",
		})
		return
	}

	if len(body.Name) < 3 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "o nome do artista deve ter ao menos 3 caracteres",
		})
		return
	}

	if len(body.Name) > 100 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "o nome do artista deve ter no máximo 100 caracteres",
		})
		return
	}

	ctx := context.Background()
	conn, err := h.db.Conn(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusConflict),
			"message": "este registro já existe",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(artist)
}

func (h *handler) ListArtist(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	conn, err := h.db.Conn(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao se conectar ao banco de dados",
		})
		return
	}

	rows, err := conn.QueryContext(ctx, `
		SELECT * FROM "artist"
	`)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao buscar os registros no banco de dados",
		})
	}

	var id int32
	var name string
	artists := []Artist{}

	for rows.Next() {
		defer rows.Close()

		if err := rows.Scan(&id, &name); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   http.StatusText(http.StatusInternalServerError),
				"message": "erro ao buscar os registros no banco de dados",
			})
			break
		}

		artists = append(artists, Artist{ID: id, Name: name})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artists)
}
