package author

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/response"
	internalValidator "go-boilerplate-rest-api-chi/internal/validator"
)

type AuthorSuccessResponse struct {
	Status  string              `json:"status" example:"success"`
	Message string              `json:"message" example:"Author retrieved successfully"`
	Author  *dto.AuthorResponse `json:"author"`
}

type AuthorHandler struct {
	service   AuthorService
	validator *internalValidator.Validator
	logger    zerolog.Logger
}

func NewAuthorHandler(service AuthorService, validator *internalValidator.Validator, logger zerolog.Logger) *AuthorHandler {
	return &AuthorHandler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *AuthorHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// routes
	r.Post("/", h.CreateAuthor)
	r.Get("/{author_id}", h.GetAuthorByID)

	return r
}

// CreateAuthor godoc
//
//	@Summary		Create a new author
//	@Description	Create a new author with the provided data
//	@Tags			authors
//	@Accept			json
//	@Produce		json
//	@Param			author	body		dto.CreateAuthorRequest	true	"Author data"
//	@Success		201		{object}	AuthorSuccessResponse
//	@Failure		400		{object}	response.ValidationErrorResponse
//	@Failure		409		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/authors [post]
func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAuthorRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		validationErrors := h.validator.FormatErrors(err)
		response.ValidationError(w, validationErrors)
		return
	}

	author, err := h.service.CreateAuthor(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, AuthorSuccessResponse{
		Status:  "success",
		Message: "Author created successfully",
		Author:  dto.ToAuthorResponse(author),
	})
}

// GetAuthorByID godoc
//
//	@Summary		Get author by id
//	@Description	Get a single author by its ID
//	@Tags			authors
//	@Produce		json
//	@Param			author_id	path		string	true	"Author ID"
//	@Success		200			{object}	AuthorSuccessResponse
//	@Failure		400			{object}	response.ErrorResponse
//	@Failure		404			{object}	response.ErrorResponse
//	@Failure		500			{object}	response.ErrorResponse
//	@Router			/authors/{author_id} [get]
func (h *AuthorHandler) GetAuthorByID(w http.ResponseWriter, r *http.Request) {
	authorID, err := uuid.Parse(chi.URLParam(r, "author_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	author, err := h.service.GetAuthorByID(r.Context(), authorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, AuthorSuccessResponse{
		Status:  "success",
		Message: "Author retrieved successfully",
		Author:  dto.ToAuthorResponse(author),
	})
}

func (h *AuthorHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, "Author not found")
	case errors.Is(err, ErrDuplicate):
		response.Error(w, http.StatusConflict, "Author with this name already exists")
	default:
		h.logger.Error().Err(err).Msg("unexpected error")
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
