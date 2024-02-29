package e2e

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/internal/style"
)

func TestStyleList(t *testing.T) {
	db_url := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	handler := style.Handler(db)

	t.Run("deve ser possível listar todos os estilos", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/style", nil)

		handler.ListStyle(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("codigo de status esperado: %d, recebido: %d", http.StatusOK, res.Code)
		}
	})
}
