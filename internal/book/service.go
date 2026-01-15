package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/book/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
)

//go:generate mockgen -destination=../mocks/mock_book_service.go -package=mocks go-boilerplate-rest-api-chi/internal/book BookService
type BookService interface {
	CreateBook(ctx context.Context, req *dto.CreateBookRequest) (*entity.Book, error)
	GetAllBooks(ctx context.Context) ([]*entity.Book, error)
	GetBookByID(ctx context.Context, bookID uuid.UUID) (*entity.Book, error)
	UpdateBook(ctx context.Context, req *dto.UpdateBookRequest, bookID uuid.UUID) error
	DeleteBook(ctx context.Context, bookID uuid.UUID) error
}

type bookService struct {
	repository       BookRepository
	authorRepository author.AuthorRepository
	logger           zerolog.Logger
}

func NewBookService(repository BookRepository, authorRepository author.AuthorRepository, logger zerolog.Logger) BookService {
	return &bookService{
		repository:       repository,
		authorRepository: authorRepository,
		logger:           logger,
	}
}

func (s *bookService) CreateBook(ctx context.Context, req *dto.CreateBookRequest) (*entity.Book, error) {
	authorID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return nil, ErrInvalidAuthorId
	}

	exists, err := s.authorRepository.Exists(ctx, authorID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, author.ErrNotFound
	}

	book := &entity.Book{
		Title:       req.Title,
		Description: req.Description,
		AuthorID:    authorID,
	}

	return s.repository.Create(ctx, book)
}

func (s *bookService) GetAllBooks(ctx context.Context) ([]*entity.Book, error) {
	books, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, ErrNotFound
	}
	return books, nil
}

func (s *bookService) GetBookByID(ctx context.Context, bookID uuid.UUID) (*entity.Book, error) {
	book, err := s.repository.GetByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) UpdateBook(ctx context.Context, req *dto.UpdateBookRequest, bookID uuid.UUID) error {
	updates := map[string]interface{}{
		"description": req.Description,
	}

	return s.repository.Update(ctx, bookID, updates)
}

func (s *bookService) DeleteBook(ctx context.Context, bookID uuid.UUID) error {
	return s.repository.Delete(ctx, bookID)
}
