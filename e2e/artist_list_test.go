package e2e

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/internal/artist"
)

func TestArtistList(t *testing.T) {
	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

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
