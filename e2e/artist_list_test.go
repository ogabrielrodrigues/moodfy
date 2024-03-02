package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/e2e/util"
	"github.com/ogabrielrodrugues/moodfy/internal/artist"
)

func TestArtistList(t *testing.T) {
	db := util.TestDatabase()
	defer util.ClearDatabase(db, "artist")

	handler := artist.Handler(db)

	t.Run("deve ser possível listar todos os artistas", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/artist", nil)

		handler.ListArtist(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusOK, res.Code)
		}
	})
}
