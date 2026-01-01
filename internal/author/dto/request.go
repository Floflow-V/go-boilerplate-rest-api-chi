package dto

type CreateAuthorRequest struct {
	Name string `json:"name" validate:"required"`
}
