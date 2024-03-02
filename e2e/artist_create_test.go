package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/e2e/util"
	"github.com/ogabrielrodrugues/moodfy/internal/artist"
)

func TestArtistCreate(t *testing.T) {
	db := util.TestDatabase()
	defer util.ClearDatabase(db, "artist")

	handler := artist.Handler(db)
	name := fmt.Sprintf("Artista Teste %d", rand.Int31n(100))

	t.Run("deve ser possível criar um artista", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]string{
			"name": name,
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/artist", &body)

		handler.CreateArtist(res, req)

		if res.Code != http.StatusCreated {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusCreated, res.Code)
		}
	})

	t.Run("não deve ser possível criar um artista se o corpo da requisição for malformado ou inválido", func(t *testing.T) {
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

	t.Run("não deve ser possível criar um artista se o nome do artista possuir menos que 3 caracteres", func(t *testing.T) {
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

	t.Run("não deve ser possível criar um artista se o nome do artista possuir mais que 100 caracteres", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]string{
			"name": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin at purus mollis, dictum felis eget, maximus ligula. Aenean efficitur facilisis.",
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/artist", &body)

		handler.CreateArtist(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possível criar um artista se ele já existir", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]string{
			"name": name,
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/artist", &body)

		handler.CreateArtist(res, req)

		if res.Code != http.StatusConflict {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusConflict, res.Code)
		}
	})
}
