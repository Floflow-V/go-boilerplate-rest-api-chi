package book

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"go-boilerplate-rest-api-chi/internal/entity"
)

//go:generate mockgen -destination=../mocks/mock_book_repository.go -package=mocks go-boilerplate-rest-api-chi/internal/book BookRepository
type BookRepository interface {
	Create(ctx context.Context, book *entity.Book) (*entity.Book, error)
	GetAll(ctx context.Context) ([]*entity.Book, error)
	GetByID(ctx context.Context, bookID uuid.UUID) (*entity.Book, error)
	Update(ctx context.Context, book *entity.Book) (*entity.Book, error)
	Delete(ctx context.Context, bookID uuid.UUID) error
}

type bookRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewBookRepository(db *gorm.DB, logger zerolog.Logger) BookRepository {
	return &bookRepository{
		db:     db,
		logger: logger,
	}
}

func (r *bookRepository) Create(ctx context.Context, newBook *entity.Book) (*entity.Book, error) {
	if err := r.db.WithContext(ctx).Create(newBook).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.logger.Error().Err(err).Msg("record already exist in database")
			return nil, ErrDuplicate
		}

		r.logger.Error().Err(err).Msg("database error")
		return nil, err
	}

	return newBook, nil
}

func (r *bookRepository) GetAll(ctx context.Context) ([]*entity.Book, error) {
	var books []*entity.Book

	if err := r.db.WithContext(ctx).Preload("Author").Find(&books).Error; err != nil {
		r.logger.Error().Err(err).Msg("error when retreive books on database ")
		return nil, err
	}

	return books, nil
}

func (r *bookRepository) GetByID(ctx context.Context, bookID uuid.UUID) (*entity.Book, error) {
	var book *entity.Book

	if err := r.db.WithContext(ctx).Preload("Author").First(&book, "id = ?", bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) Update(ctx context.Context, book *entity.Book) (*entity.Book, error) {
	if err := r.db.WithContext(ctx).Save(book).Error; err != nil {
		return nil, err
	}

	return book, nil
}

func (r *bookRepository) Delete(ctx context.Context, bookID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", bookID).Delete(&entity.Book{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
