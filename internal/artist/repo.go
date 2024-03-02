package artist

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
)

type repo struct {
	db *sql.DB
}

func Repo(db *sql.DB) *repo {
	return &repo{db}
}

func (r *repo) Create(body *DTO) (*Artist, int, error) {
	if len(body.Name) < 3 {
		return nil, http.StatusUnprocessableEntity, errors.New("o nome do artista deve ter ao menos 3 caracteres")
	}

	if len(body.Name) > 100 {
		return nil, http.StatusUnprocessableEntity, errors.New("o nome do artista deve ter no máximo 100 caracteres")
	}

	ctx := context.Background()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao se conectar ao banco de dados")
	}

	artist := New(body.Name)

	row := conn.QueryRowContext(ctx, `
		INSERT INTO "artist" (name)
		VALUES ($1) RETURNING id`,
		artist.Name,
	)

	if err := row.Scan(&artist.ID); err != nil {
		return nil, http.StatusConflict, errors.New("este registro já existe")
	}

	return artist, http.StatusCreated, nil
}

func (r *repo) List() (*[]Artist, int, error) {
	ctx := context.Background()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao se conectar ao banco de dados")
	}

	rows, err := conn.QueryContext(ctx, `
		SELECT * FROM "artist"
	`)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao buscar os registros no banco de dados")
	}

	var id int32
	var name string
	artists := []Artist{}

	for rows.Next() {
		defer rows.Close()

		if err := rows.Scan(&id, &name); err != nil {
			return nil, http.StatusInternalServerError, errors.New("erro ao buscar os registros no banco de dados")
		}

		artists = append(artists, Artist{ID: id, Name: name})
	}

	return &artists, http.StatusOK, nil
}
