package dto

type CreateBookRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	AuthorID    string `json:"author_id" validate:"required"`
}

type UpdateBookRequest struct {
	Description string `json:"description" validate:"required"`
}
