package music

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func string_to_arr(str *string) []string {
	if *str == "{}" {
		return nil
	}

	*str = strings.TrimLeft(*str, "{")
	*str = strings.TrimRight(*str, "}")
	*str = strings.Join(strings.Split(*str, `"`), "")

	return strings.Split(*str, ",")
}

type repo struct {
	db *sql.DB
}

func Repo(db *sql.DB) *repo {
	return &repo{db}
}

func (r *repo) Create(body *DTO) (*Music, int, error) {
	if len(body.Name) < 3 {
		return nil, http.StatusUnprocessableEntity, errors.New("o nome da música deve ter ao menos 3 caracteres")
	}

	if len(body.Name) > 100 {
		return nil, http.StatusUnprocessableEntity, errors.New("o nome da música deve ter no máximo 50 caracteres")
	}

	if _, err := url.Parse(body.CoverImage); err != nil {
		return nil, http.StatusUnprocessableEntity, errors.New("a capa da música deve ser uma url válida")
	}

	if len(body.CoverImage) > 300 {
		return nil, http.StatusUnprocessableEntity, errors.New("a url da capa da música deve ter no máximo 300 caracteres")
	}

	if _, err := url.Parse(body.SpotifyLink); err != nil {
		return nil, http.StatusUnprocessableEntity, errors.New("a url da música no spotify deve ser uma url válida")
	}

	if len(body.SpotifyLink) > 200 {
		return nil, http.StatusUnprocessableEntity, errors.New("a url da música no spotify deve ter no máximo 200 caracteres")
	}

	ctx := context.Background()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao se conectar ao banco de dados")
	}

	music := New(body.Name, body.CoverImage, body.SpotifyLink)

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return nil, http.StatusInternalServerError, errors.New("erro ao iniciar a transação com o banco de dados")
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO "music" (artist_id, name, link, cover)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		body.ArtistID,
		music.Name,
		music.SpotifyLink,
		music.CoverImage,
	)

	if err := row.Scan(&music.ID); err != nil {
		tx.Rollback()
		return nil, http.StatusConflict, errors.New("este registro já existe")
	}

	for _, style_id := range body.Styles {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO "music_style" (music_id, style_id)
			VALUES ($1, $2)`,
			music.ID,
			style_id,
		); err != nil {
			tx.Rollback()
			return nil, http.StatusInternalServerError, errors.New("erro ao inserir no banco de dados")
		}
	}

	tx.Commit()
	return music, http.StatusCreated, nil
}

func (r *repo) List(filters *url.Values) (*[]Music, int, error) {
	ctx := context.Background()
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao se conectar ao banco de dados")
	}

	var query string
	query = fmt.Sprintf(`
		SELECT
			m.id, 
			m.name AS music_name, 
			a.name AS artist_name, 
			m.cover AS cover_image, 
			m.link AS spotify_link, 
			ARRAY_AGG(s.name) AS styles 
		FROM "music" m 
		JOIN "artist" a 
		ON m.artist_id = a.id 
		JOIN "music_style" ms 
		ON m.id = ms.music_id 
		JOIN "style" s 
		ON ms.style_id = s.id`)

	if filters.Has("artist") && filters.Has("style") {
		query = fmt.Sprintf(`%s 
		WHERE a.name = '%s' 
		AND s.name = '%s'`, query, filters.Get("artist"), filters.Get("style"))
	} else {
		if filters.Has("artist") {
			query = fmt.Sprintf(`%s 
		WHERE a.name = '%s'`, query, filters.Get("artist"))
		}

		if filters.Has("style") {
			query = fmt.Sprintf(`%s 
		WHERE s.name = '%s'`, query, filters.Get("style"))
		}
	}

	query = fmt.Sprintf(`
		%s %s
	`, query, "GROUP BY m.id, a.name")

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("erro ao buscar os registros no banco de dados")
	}

	var id int32
	var music_name, artist_name, cover_image, spotify_link string
	var styles *string
	musics := []Music{}

	for rows.Next() {
		defer rows.Close()

		if err := rows.Scan(
			&id,
			&music_name,
			&artist_name,
			&cover_image,
			&spotify_link,
			&styles,
		); err != nil {
			return nil, http.StatusInternalServerError, errors.New("erro ao buscar os registros no banco de dados")
		}

		musics = append(musics, Music{
			ID:          id,
			Artist:      artist_name,
			Name:        music_name,
			CoverImage:  cover_image,
			SpotifyLink: spotify_link,
			Styles:      string_to_arr(styles),
		})
	}

	return &musics, http.StatusOK, nil
}
