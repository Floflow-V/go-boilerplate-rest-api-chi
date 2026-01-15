package dto

import (
	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
)

type BookResponse struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Author      dto.AuthorResponse `json:"author,omitempty"`
}

func ToBookResponse(book *entity.Book) *BookResponse {
	return &BookResponse{
		ID:          book.ID.String(),
		Title:       book.Title,
		Description: book.Description,
		Author: dto.AuthorResponse{
			ID:   book.Author.ID.String(),
			Name: book.Author.Name,
		},
	}
}

func ToBooksResponse(books []*entity.Book) []BookResponse {
	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *ToBookResponse(book)
	}
	return responses
}
