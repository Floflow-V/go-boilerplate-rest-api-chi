package dto

import "go-boilerplate-rest-api-chi/internal/entity"

type AuthorResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToAuthorResponse(author *entity.Author) *AuthorResponse {
	return &AuthorResponse{
		ID:   author.ID.String(),
		Name: author.Name,
	}
}
