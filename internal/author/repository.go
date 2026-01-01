package author

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"go-boilerplate-rest-api-chi/internal/entity"
)

//go:generate mockgen -destination=../mocks/mock_author_repository.go -package=mocks go-boilerplate-rest-api-chi/internal/author AuthorRepository
type AuthorRepository interface {
	Create(ctx context.Context, newAuthor *entity.Author) (*entity.Author, error)
	GetByID(ctx context.Context, authorID uuid.UUID) (*entity.Author, error)
}

type authorRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewAuthorRepository(db *gorm.DB, logger zerolog.Logger) AuthorRepository {
	return &authorRepository{
		db:     db,
		logger: logger,
	}
}

func (r *authorRepository) Create(ctx context.Context, newAuthor *entity.Author) (*entity.Author, error) {
	if err := r.db.WithContext(ctx).Create(newAuthor).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrDuplicate
		}

		r.logger.Error().Err(err).Msg("database error")
		return nil, err
	}

	return newAuthor, nil
}

func (r *authorRepository) GetByID(ctx context.Context, authorID uuid.UUID) (*entity.Author, error) {
	var author *entity.Author

	if err := r.db.WithContext(ctx).First(&author, "id = ?", authorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}

		r.logger.Error().Err(err).Msg("database error")
		return nil, err
	}

	return author, nil
}
