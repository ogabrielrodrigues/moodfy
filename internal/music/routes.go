package music

import (
	"context"
	"database/sql"
	"encoding/json"
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

type handler struct {
	db *sql.DB
}

func Handler(db *sql.DB) *handler {
	return &handler{db}
}

func (h *handler) CreateMusic(w http.ResponseWriter, r *http.Request) {
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
			"message": "o nome da música deve ter ao menos 3 caracteres",
		})
		return
	}

	if len(body.Name) > 100 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "o nome da música deve ter no máximo 50 caracteres",
		})
		return
	}

	if _, err := url.Parse(body.CoverImage); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "a capa da música deve ser uma url válida",
		})
		return
	}

	if len(body.CoverImage) > 300 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "a url da capa da música deve ter no máximo 300 caracteres",
		})
		return
	}

	if _, err := url.Parse(body.SpotifyLink); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "a url da música no spotify deve ser uma url válida",
		})
		return
	}

	if len(body.SpotifyLink) > 200 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusUnprocessableEntity),
			"message": "a url da música no spotify deve ter no máximo 200 caracteres",
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

	music := New(body.ArtistID, body.Name, body.CoverImage, body.SpotifyLink)

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao iniciar a transação com o  banco de dados",
		})
		return
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO "music" (artist_id, name, link, cover)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		music.ArtistID,
		music.Name,
		music.SpotifyLink,
		music.CoverImage,
	)

	if err := row.Scan(&music.ID); err != nil {
		fmt.Println(err)
		tx.Rollback()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusConflict),
			"message": "este registro já existe",
		})
		return
	}

	for _, style_id := range body.Styles {
		if _, err := tx.ExecContext(ctx, `
		INSERT INTO "music_style" (music_id, style_id)
		VALUES ($1, $2)`,
			music.ID,
			style_id,
		); err != nil {
			tx.Rollback()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   http.StatusText(http.StatusInternalServerError),
				"message": "erro ao inserir no banco de dados",
			})
			return
		}
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(music)
}

// List or filter musics
func (h *handler) ListMusic(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

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

	if params.Has("artist") && params.Has("style") {
		query = fmt.Sprintf(`%s 
		WHERE a.name = '%s' 
		AND s.name = '%s'`, query, params.Get("artist"), params.Get("style"))
	} else {
		if params.Has("artist") {
			query = fmt.Sprintf(`%s 
		WHERE a.name = '%s'`, query, params.Get("artist"))
		}

		if params.Has("style") {
			query = fmt.Sprintf(`%s 
		WHERE s.name = '%s'`, query, params.Get("style"))
		}
	}

	query = fmt.Sprintf(`
		%s %s
	`, query, "GROUP BY m.id, a.name")

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   http.StatusText(http.StatusInternalServerError),
			"message": "erro ao buscar os registros no banco de dados",
		})
	}

	var id int32
	var music_name, artist_name, cover_image, spotify_link string
	var styles *string
	musics := []map[string]any{}
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   http.StatusText(http.StatusInternalServerError),
				"message": "erro ao buscar os registros no banco de dados",
			})
			break
		}

		musics = append(musics, map[string]any{
			"id":           id,
			"artist":       artist_name,
			"music":        music_name,
			"cover_image":  cover_image,
			"spotify_link": spotify_link,
			"styles":       string_to_arr(styles),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(musics)
}
