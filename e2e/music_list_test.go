package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ogabrielrodrugues/moodfy/e2e/util"
	"github.com/ogabrielrodrugues/moodfy/internal/music"
)

func TestMusicList(t *testing.T) {
	db := util.TestDatabase()
	defer util.ClearDatabase(db, "music")

	handler := music.Handler(db)

	t.Run("deve ser possível listar todas as músicas", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/music", nil)

		handler.ListMusic(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusOK, res.Code)
		}
	})

	t.Run("deve ser possível filtrar as músicas por artista", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/music?artist=zezin", nil)

		handler.ListMusic(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusOK, res.Code)
		}
	})

	t.Run("deve ser possível filtrar as músicas por estilo", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/music?style=qualquer", nil)

		handler.ListMusic(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusOK, res.Code)
		}
	})

	t.Run("deve ser possível filtrar as músicas por artista e estilo", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/music?artist=zezin&style=qualquer", nil)

		handler.ListMusic(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusOK, res.Code)
		}
	})
}
