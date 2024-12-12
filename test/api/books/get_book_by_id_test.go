package books

import (
	app "book_service/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBookById(t *testing.T) {
	r := app.SetupServer()

	t.Run("Valid Book ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/books/valid-uuid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("Invalid Book ID Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/books/invalid-id", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("Non-existent Book ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/books/non-existent-uuid", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
	})
}
