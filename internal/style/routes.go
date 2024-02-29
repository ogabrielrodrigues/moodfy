package style

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

func (h *handler) CreateStyle(w http.ResponseWriter, r *http.Request) {
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
			"message": "o nome do estilo deve ter ao menos 3 caracteres",
		})
		return
	}

	if len(body.Name) > 50 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "o nome do estilo deve ter no máximo 50 caracteres",
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

	style := New(body.Name)

	row := conn.QueryRowContext(ctx, `
		INSERT INTO "style" (name)
		VALUES ($1) RETURNING id`,
		style.Name,
	)

	if err := row.Scan(&style.ID); err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusConflict),
			"message": "este registro já existe",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(style)
}

func (h *handler) ListStyle(w http.ResponseWriter, r *http.Request) {
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

	rows, err := conn.QueryContext(ctx, `
		SELECT * FROM "style"
	`)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao buscar os registros no banco de dados",
		})
	}

	var id int32
	var name string
	styles := []Style{}

	for rows.Next() {
		defer rows.Close()

		if err := rows.Scan(&id, &name); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   http.StatusText(http.StatusInternalServerError),
				"message": "erro ao buscar os registros no banco de dados",
			})
			break
		}

		styles = append(styles, Style{ID: id, Name: name})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(styles)
}
