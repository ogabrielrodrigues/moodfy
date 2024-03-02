package style

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

func (r *repo) Create(body *DTO) (*Style, int, error) {
	if len(body.Name) < 3 {
		return nil, http.StatusUnprocessableEntity, errors.New("o nome do estilo deve ter ao menos 3 caracteres")
	}

	if len(body.Name) > 50 {
		return nil, http.StatusUnprocessableEntity, errors.New("o nome do estilo deve ter no máximo 50 caracteres")
	}

	ctx := context.Background()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao se conectar ao banco de dados")
	}

	style := New(body.Name)

	row := conn.QueryRowContext(ctx, `
		INSERT INTO "style" (name)
		VALUES ($1) RETURNING id`,
		style.Name,
	)

	if err := row.Scan(&style.ID); err != nil {
		return nil, http.StatusConflict, errors.New("este registro já existe")
	}

	return style, http.StatusCreated, nil
}

func (r *repo) List() (*[]Style, int, error) {
	ctx := context.Background()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao se conectar ao banco de dados")
	}

	rows, err := conn.QueryContext(ctx, `
		SELECT * FROM "style"
	`)

	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao buscar os registros no banco de dados")
	}

	var id int32
	var name string
	styles := []Style{}

	for rows.Next() {
		defer rows.Close()

		if err := rows.Scan(&id, &name); err != nil {
			return nil, http.StatusInternalServerError, errors.New("erro ao buscar os registros no banco de dados")
		}

		styles = append(styles, Style{ID: id, Name: name})
	}

	return &styles, http.StatusOK, nil
}
