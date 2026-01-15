package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-boilerplate-rest-api-chi/internal/api"
	"go-boilerplate-rest-api-chi/internal/config"
	"go-boilerplate-rest-api-chi/internal/entity"
)

func TestCreateApi(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db.AutoMigrate(&entity.Book{}, &entity.Author{})

	t.Run("development_mode", func(t *testing.T) {
		cfg := config.Config{Api: config.ApiConfig{Environment: "development"}}
		handler := api.CreateApi(cfg, zerolog.Nop(), db)

		req := httptest.NewRequest(http.MethodGet, "/api/alive", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		req = httptest.NewRequest(http.MethodGet, "/api/doc/index.html", nil)
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		req = httptest.NewRequest(http.MethodOptions, "/api/books", nil)
		req.Header.Set("Origin", "http://localhost")
		req.Header.Set("Access-Control-Request-Method", "GET")
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.NotEmpty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("production_mode", func(t *testing.T) {
		cfg := config.Config{Api: config.ApiConfig{Environment: "production"}}
		handler := api.CreateApi(cfg, zerolog.Nop(), db)

		req := httptest.NewRequest(http.MethodGet, "/api/doc/index.html", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
