package book

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/book/dto"
	"go-boilerplate-rest-api-chi/internal/response"
	internalValidator "go-boilerplate-rest-api-chi/internal/validator"
)

type BookSuccessResponse struct {
	Status  string            `json:"status" example:"success"`
	Message string            `json:"message" example:"Book retrieved successfully"`
	Book    *dto.BookResponse `json:"book"`
}

type BooksSuccessResponse struct {
	Status  string             `json:"status" example:"success"`
	Message string             `json:"message" example:"Books retrieved successfully"`
	Books   []dto.BookResponse `json:"books"`
}

type BookHandler struct {
	service   BookService
	validator *internalValidator.Validator
	logger    zerolog.Logger
}

func NewBookHandler(service BookService, validator *internalValidator.Validator, logger zerolog.Logger) *BookHandler {
	return &BookHandler{
		service:   service,
		validator: validator,
		logger:    logger,
	}
}

func (h *BookHandler) Routes() http.Handler {
	r := chi.NewRouter()

	// routes
	r.Post("/", h.CreateBook)
	r.Get("/", h.GetAllBooks)
	r.Get("/{book_id}", h.GetBookByID)
	r.Put("/{book_id}", h.UpdateBook)
	r.Get("/secure", h.AuthTestRoute)

	return r
}

// CreateBook godoc
//
//	@Summary		Create a new book
//	@Description	Create a new book with the provided data
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book	body		dto.CreateBookRequest	true	"Book data"
//	@Success		201		{object}	BookSuccessResponse
//	@Failure		400		{object}	response.ValidationErrorResponse
//	@Failure		409		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/books [post]
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		validationErrors := h.validator.FormatErrors(err)
		response.ValidationError(w, validationErrors)
		return
	}

	book, err := h.service.CreateBook(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, BookSuccessResponse{
		Status:  "success",
		Message: "Book created successfully",
		Book:    dto.ToBookResponse(book),
	})
}

// GetAllBooks godoc
//
//	@Summary		Get all books
//	@Description	Get a list of all books
//	@Tags			books
//	@Produce		json
//	@Success		200	{object}	BooksSuccessResponse
//	@Failure		404	{object}	response.ErrorResponse
//	@Failure		500	{object}	response.ErrorResponse
//	@Router			/books [get]
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAllBooks(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, BooksSuccessResponse{
		Status:  "success",
		Message: "Books retrieved successfully",
		Books:   dto.ToBooksResponse(books),
	})
}

// GetBookByID godoc
//
//	@Summary		Get book by id
//	@Description	Get a single book by its ID
//	@Tags			books
//	@Produce		json
//	@Param			book_id	path		string	true	"Book ID"
//	@Success		200		{object}	BookSuccessResponse
//	@Failure		400		{object}	response.ErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/books/{book_id} [get]
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	bookID, err := uuid.Parse(chi.URLParam(r, "book_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	book, err := h.service.GetBookByID(r.Context(), bookID)
	if err != nil {
		h.handleError(w, err)
	}

	response.JSON(w, http.StatusOK, BookSuccessResponse{
		Status:  "success",
		Message: "Book retrieved successfully",
		Book:    dto.ToBookResponse(book),
	})
}

// UpdateBook godoc
//
//	@Summary		Update a book
//	@Description	Update a book with the provided data
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			book_id	path		string					true	"Book ID"
//	@Param			book	body		dto.UpdateBookRequest	true	"Book data"
//	@Success		200		{object}	BookSuccessResponse
//	@Failure		400		{object}	response.ValidationErrorResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/books/{book_id} [put]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := uuid.Parse(chi.URLParam(r, "book_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	var req dto.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		validationErrors := h.validator.FormatErrors(err)
		response.ValidationError(w, validationErrors)
		return
	}

	book, err := h.service.UpdateBook(r.Context(), &req, bookID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, BookSuccessResponse{
		Status:  "success",
		Message: "Book updated successfully",
		Book:    dto.ToBookResponse(book),
	})
}

// DeleteBook godoc
//
//	@Summary		Delete a book
//	@Description	Delete a book by its ID
//	@Tags			books
//	@Produce		json
//	@Param			book_id	path		string	true	"Book ID"
//	@Success		200		{object}	response.SuccessResponse
//	@Failure		404		{object}	response.ErrorResponse
//	@Failure		500		{object}	response.ErrorResponse
//	@Router			/books/{book_id} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID, err := uuid.Parse(chi.URLParam(r, "book_id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid uuid")
		return
	}

	err = h.service.DeleteBook(r.Context(), bookID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, "Book deleted successfully")
}

// AuthTestRoute godoc
//
//	@Summary		Authenticated test route
//	@Description	Authenticated test route
//	@Tags			books
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	response.SuccessResponse
//	@Router			/secure [get]
func (h *BookHandler) AuthTestRoute(w http.ResponseWriter, r *http.Request) {
	response.Success(w, "ok")
}

func (h *BookHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, "Book not found")
	case errors.Is(err, ErrDuplicate):
		response.Error(w, http.StatusConflict, "Book with this name already exists")
	case errors.Is(err, ErrInvalidAuthorId):
		response.Error(w, http.StatusBadRequest, "invalid author ID")
	case errors.Is(err, author.ErrNotFound):
		response.Error(w, http.StatusNotFound, "Author not found")
	default:
		h.logger.Error().Err(err).Msg("unexpected error")
		response.Error(w, http.StatusInternalServerError, "Internal server error")
	}
}
