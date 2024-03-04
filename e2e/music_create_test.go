package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogabrielrodrugues/moodfy/e2e/util"
	"github.com/ogabrielrodrugues/moodfy/internal/artist"
	"github.com/ogabrielrodrugues/moodfy/internal/music"
	"github.com/ogabrielrodrugues/moodfy/internal/style"
)

func TestMusicCreate(t *testing.T) {
	db := util.TestDatabase()
	defer util.ClearDatabase(db, "music")
	defer util.ClearDatabase(db, "style")
	defer util.ClearDatabase(db, "artist")

	artist, _, _ := artist.Repo(db).Create(&artist.DTO{
		Name: fmt.Sprintf("Artista %d", rand.Int31n(100)),
	})

	style1, _, _ := style.Repo(db).Create(&style.DTO{
		Name: fmt.Sprintf("Estilo %d", rand.Int31n(100)),
	})

	style2, _, _ := style.Repo(db).Create(&style.DTO{
		Name: fmt.Sprintf("Estilo %d", rand.Int31n(100)),
	})

	style3, _, _ := style.Repo(db).Create(&style.DTO{
		Name: fmt.Sprintf("Estilo %d", rand.Int31n(100)),
	})

	handler := music.Handler(db)
	name := fmt.Sprintf("Música %d", rand.Int31n(100))
	t.Run("deve ser possivel criar uma música", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         name,
			"cover_image":  fmt.Sprintf("http://coverimages.com/%d.png", rand.Int31n(100)),
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusCreated {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusCreated, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se o corpo da requisição for malformado ou inválido", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    "ID String",
			"name":         fmt.Sprintf("Música %d", rand.Int31n(100)),
			"cover_image":  fmt.Sprintf("http://coverimages.com/%d.png", rand.Int31n(100)),
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusBadRequest, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se o nome da música possuir menos de 3 caracteres", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         "dd",
			"cover_image":  fmt.Sprintf("http://coverimages.com/%d.png", rand.Int31n(100)),
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se o nome da música possuir mais que 50 caracteres", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin at purus mollis, dictum felis eget, maximus ligula. Aenean efficitur facilisis.",
			"cover_image":  fmt.Sprintf("http://coverimages.com/%d.png", rand.Int31n(100)),
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se a url da capa for inválida", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         fmt.Sprintf("Música %d", rand.Int31n(100)),
			"cover_image":  "dsdsds",
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se a url da capa possuir mais que 300 caracteres", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         fmt.Sprintf("Música %d", rand.Int31n(100)),
			"cover_image":  "https://crawler-test.com/urls/page_url_length/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se a url do spotify for inválida", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         fmt.Sprintf("Música %d", rand.Int31n(100)),
			"cover_image":  fmt.Sprintf("http://coverimages.com/%d.png", rand.Int31n(100)),
			"spotify_link": "dsdsds",
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se a url do spotify possuir mais que 200 caracteres", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         fmt.Sprintf("Música %d", rand.Int31n(100)),
			"cover_image":  "https://crawler-test.com/urls/page_url_length/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusUnprocessableEntity {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusUnprocessableEntity, res.Code)
		}
	})

	t.Run("não deve ser possivel criar uma música se ela já existir", func(t *testing.T) {
		var body bytes.Buffer

		json.NewEncoder(&body).Encode(map[string]any{
			"artist_id":    artist.ID,
			"name":         name,
			"cover_image":  fmt.Sprintf("http://coverimages.com/%d.png", rand.Int31n(100)),
			"spotify_link": fmt.Sprintf("http://musiclinks.com/%d", rand.Int31n(100)),
			"styles": []int32{
				style1.ID,
				style2.ID,
				style3.ID,
			},
		})

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/music", &body)

		handler.CreateMusic(res, req)

		if res.Code != http.StatusConflict {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusConflict, res.Code)
		}
	})
}
