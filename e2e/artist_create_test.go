package e2e

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/internal/artist"
)

func TestArtistCreate(t *testing.T) {
	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	handler := artist.Handler(db)

	t.Run("deve ser possivel criar um artista", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]string{
			"name": fmt.Sprintf("Artista Qualquer %d", rand.Int31n(100)),
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/artist", &body)

		handler.CreateArtist(res, req)

		if res.Code != http.StatusCreated {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusCreated, res.Code)
		}
	})

	t.Run("não deve ser possivel criar um artista se o corpo da requisição for malformado ou inválido", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]int{
			"name": 43,
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/artist", &body)

		handler.CreateArtist(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusBadRequest, res.Code)
		}
	})

	t.Run("não deve ser possivel criar um artista se o nome do artista possuir menos que 3 caracteres", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]string{
			"name": "Ar",
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/artist", &body)

		handler.CreateArtist(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})
}
