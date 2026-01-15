package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	authorDTO "go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/book/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
)

func TestToAuthorResponse(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		entity := entity.Book{
			ID:          uuid.MustParse("58411bf8-aa11-4553-9b13-4bdf58875d35"),
			Title:       "Book1",
			Description: "Description1",
			AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
			Author: &entity.Author{
				ID:        uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
				Name:      "Author1",
				CreatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
				UpdatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
			},
			CreatedAt: time.Date(2026, 01, 12, 21, 45, 00, 00, time.Local),
			UpdatedAt: time.Date(2026, 01, 12, 21, 45, 00, 00, time.Local),
		}

		expectedResponse := dto.BookResponse{
			ID:          "58411bf8-aa11-4553-9b13-4bdf58875d35",
			Title:       "Book1",
			Description: "Description1",
			Author: authorDTO.AuthorResponse{
				ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
				Name: "Author1",
			},
		}

		response := dto.ToBookResponse(&entity)

		assert.Equal(t, &expectedResponse, response)
	})
}

func TestToBooksResponse(t *testing.T) {
	t.Run("nominal for multiple books", func(t *testing.T) {
		entities := []*entity.Book{
			{
				ID:          uuid.MustParse("58411bf8-aa11-4553-9b13-4bdf58875d35"),
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
				Author: &entity.Author{
					ID:        uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
					Name:      "Author1",
					CreatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
					UpdatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
				},
				CreatedAt: time.Date(2026, 01, 12, 21, 45, 00, 00, time.Local),
				UpdatedAt: time.Date(2026, 01, 12, 21, 45, 00, 00, time.Local),
			},
			{
				ID:          uuid.MustParse("26dcbb76-e09d-45c5-8c97-f32ce3a2766a"),
				Title:       "Book2",
				Description: "Description2",
				AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
				Author: &entity.Author{
					ID:        uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
					Name:      "Author1",
					CreatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
					UpdatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
				},
				CreatedAt: time.Date(2026, 01, 12, 21, 50, 00, 00, time.Local),
				UpdatedAt: time.Date(2026, 01, 12, 21, 50, 00, 00, time.Local),
			},
			{
				ID:          uuid.MustParse("2933e943-bde9-4743-8961-97828d166e11"),
				Title:       "Book3",
				Description: "Description3",
				AuthorID:    uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
				Author: &entity.Author{
					ID:        uuid.MustParse("b846fc59-401a-450d-b3f1-3e9a953d7c22"),
					Name:      "Author1",
					CreatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
					UpdatedAt: time.Date(2026, 01, 10, 21, 45, 00, 00, time.Local),
				},
				CreatedAt: time.Date(2026, 01, 12, 21, 55, 00, 00, time.Local),
				UpdatedAt: time.Date(2026, 01, 12, 21, 55, 00, 00, time.Local),
			},
		}

		expectedResponse := []dto.BookResponse{
			{
				ID:          "58411bf8-aa11-4553-9b13-4bdf58875d35",
				Title:       "Book1",
				Description: "Description1",
				Author: authorDTO.AuthorResponse{
					ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
					Name: "Author1",
				},
			},
			{
				ID:          "26dcbb76-e09d-45c5-8c97-f32ce3a2766a",
				Title:       "Book2",
				Description: "Description2",
				Author: authorDTO.AuthorResponse{
					ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
					Name: "Author1",
				},
			},
			{
				ID:          "2933e943-bde9-4743-8961-97828d166e11",
				Title:       "Book3",
				Description: "Description3",
				Author: authorDTO.AuthorResponse{
					ID:   "b846fc59-401a-450d-b3f1-3e9a953d7c22",
					Name: "Author1",
				},
			},
		}

		response := dto.ToBooksResponse(entities)

		assert.Equal(t, expectedResponse, response)
	})
}
