package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"

	"github.com/ogabrielrodrugues/moodfy/e2e/util"
	"github.com/ogabrielrodrugues/moodfy/internal/style"
)

func TestStyleList(t *testing.T) {
	db := util.TestDatabase()
	defer util.ClearDatabase(db, "style")

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
