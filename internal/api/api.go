package api

import (
	"net/http"
	"os"
	"time"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/book"
	"go-boilerplate-rest-api-chi/internal/config"
	internalValidator "go-boilerplate-rest-api-chi/internal/validator"
)

func CreateApi(cfg config.Config, logger zerolog.Logger, db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.CleanPath,
		middleware.StripSlashes,
		middleware.GetHead,
		middleware.Timeout(10*time.Second),
		middleware.Throttle(100), // limit the number of request globaly for all the api
		httprate.LimitByRealIP(100, 1*time.Minute),
	)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           12 * int(time.Hour),
	}))

	api := chi.NewRouter()

	api.Use(middleware.Heartbeat("/api/alive"))

	validator := internalValidator.New()

	// -------- Repos / Services / Handlers --------

	bookRepo := book.NewBookRepository(db, logger)
	authorRepo := author.NewAuthorRepository(db, logger)

	bookService := book.NewBookService(bookRepo, authorRepo, logger)
	authorService := author.NewAuthorService(authorRepo, logger)

	bookHandler := book.NewBookHandler(bookService, validator, logger)
	authorHandler := author.NewAuthorHandler(authorService, validator, logger)

	api.Mount("/books", bookHandler.Routes())
	api.Mount("/authors", authorHandler.Routes())

	if cfg.Api.Environment == "development" {
		api.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
			content, err := os.ReadFile("./docs/swagger.json")
			if err != nil {
				logger.Error().Err(err).Msg("Impossible de lire le fichier swagger.json")
				http.Error(w, "Documentation non disponible", http.StatusInternalServerError)
				return
			}

			htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
				SpecContent: string(content),

				CustomOptions: scalar.CustomOptions{
					PageTitle: "Boilerplate - rest API - CHI",
				},
				DarkMode:      true,
				IsEditable:    false,
				Theme:         scalar.ThemeBluePlanet,
				HideModels:    true,
				BaseServerURL: "/api",
			})
			if err != nil {
				logger.Error().Err(err).Msg("Scalar Documentation error")
			}

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, err = w.Write([]byte(htmlContent))
			if err != nil {
				panic(err)
			}
		})
	}

	if cfg.Api.Environment == "development" {
		api.Get("/doc/*", httpSwagger.WrapHandler)
	}

	r.Mount("/api", api)

	return r
}
