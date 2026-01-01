package dto

import (
	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
)

type BookResponse struct {
	ID     string              `json:"id"`
	Title  string              `json:"title"`
	Author *dto.AuthorResponse `json:"author,omitempty"`
}

func ToBookResponse(book *entity.Book) *BookResponse {
	var author *dto.AuthorResponse

	if book.Author != nil {
		author = &dto.AuthorResponse{
			ID:   book.Author.ID.String(),
			Name: book.Author.Name,
		}
	}

	return &BookResponse{
		ID:     book.ID.String(),
		Title:  book.Title,
		Author: author,
	}
}

func ToBooksResponse(books []*entity.Book) []BookResponse {
	responses := make([]BookResponse, len(books))
	for i, book := range books {
		responses[i] = *ToBookResponse(book)
	}
	return responses
}
