package dto_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
)

func TestToAuthorResponse(t *testing.T) {
	t.Run("nominal", func(t *testing.T) {
		entity := entity.Author{
			ID:   uuid.MustParse("aeca0955-bae4-47e9-9f85-6818dc68ca51"),
			Name: "George R.R. Martin",
		}

		expectedrResponse := dto.AuthorResponse{
			ID:   "aeca0955-bae4-47e9-9f85-6818dc68ca51",
			Name: "George R.R. Martin",
		}

		response := dto.ToAuthorResponse(&entity)

		assert.Equal(t, &expectedrResponse, response)
	})
}
