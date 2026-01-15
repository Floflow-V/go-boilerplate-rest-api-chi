package author_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/entity"
	"go-boilerplate-rest-api-chi/internal/test-utils"
)

func TestAuthorRepository_Create(t *testing.T) {
	tests := []struct {
		name             string
		input            *entity.Author
		configureMock    func(sqlmock.Sqlmock, *entity.Author)
		expectedError    error
		expectedResponse *entity.Author
	}{
		{
			name: "success create author",
			input: &entity.Author{
				Name: "Victor Hugo",
			},
			configureMock: func(mock sqlmock.Sqlmock, input *entity.Author) {
				mock.ExpectExec("INSERT INTO `authors`").
					WithArgs(
						sqlmock.AnyArg(), // ID généré
						input.Name,
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
			expectedResponse: &entity.Author{
				Name: "Victor Hugo",
			},
		},
		{
			name: "error duplicate author",
			input: &entity.Author{
				Name: "Duplicate Author",
			},
			configureMock: func(mock sqlmock.Sqlmock, input *entity.Author) {
				mock.ExpectExec("INSERT INTO `authors`").
					WithArgs(
						sqlmock.AnyArg(), // ID généré
						input.Name,
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnError(gorm.ErrDuplicatedKey)
			},
			expectedError:    author.ErrDuplicate,
			expectedResponse: nil,
		},
		{
			name: "error database connection failed",
			input: &entity.Author{
				Name: "Test Author",
			},
			configureMock: func(mock sqlmock.Sqlmock, input *entity.Author) {
				mock.ExpectExec("INSERT INTO `authors`").
					WithArgs(
						sqlmock.AnyArg(), // ID généré
						input.Name,
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectedError:    gorm.ErrInvalidDB,
			expectedResponse: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock, test.input)

			repo := author.NewAuthorRepository(db, zerolog.Nop())

			newAuthor, err := repo.Create(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, newAuthor)
				assert.Equal(t, test.expectedResponse.Name, newAuthor.Name)
				assert.NotEqual(t, uuid.Nil, newAuthor.ID)
				assert.NotZero(t, newAuthor.CreatedAt)
				assert.NotZero(t, newAuthor.UpdatedAt)
			} else {
				assert.Nil(t, newAuthor)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthorRepository_GetByID(t *testing.T) {
	tests := []struct {
		name             string
		authorID         uuid.UUID
		configureMock    func(sqlmock.Sqlmock, uuid.UUID)
		expectedError    error
		expectedResponse *entity.Author
	}{
		{
			name:     "success get author by id",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				now := time.Now()

				rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(id, "Victor Hugo", now, now)

				mock.ExpectQuery("SELECT \\* FROM `authors` WHERE id = \\? ORDER BY `authors`.`id` LIMIT \\?").
					WithArgs(id, 1).
					WillReturnRows(rows)
			},
			expectedError: nil,
			expectedResponse: &entity.Author{
				ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
				Name: "Victor Hugo",
			},
		},
		{
			name:     "error author not found",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT \* FROM .authors. WHERE id = \? ORDER BY .authors.\..id. LIMIT \?`).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedError:    author.ErrNotFound,
			expectedResponse: nil,
		},
		{
			name:     "error database connection failed",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT \* FROM .authors. WHERE id = \? ORDER BY .authors.\..id. LIMIT \?`).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectedError:    gorm.ErrInvalidDB,
			expectedResponse: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock, test.authorID)

			repo := author.NewAuthorRepository(db, zerolog.Nop())

			author, err := repo.GetByID(context.Background(), test.authorID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, author)
				assert.Equal(t, test.expectedResponse.ID, author.ID)
				assert.Equal(t, test.expectedResponse.Name, author.Name)
				assert.NotZero(t, author.CreatedAt)
				assert.NotZero(t, author.UpdatedAt)
			} else {
				assert.Nil(t, author)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthorRepository_Exists(t *testing.T) {
	tests := []struct {
		name          string
		authorID      uuid.UUID
		configureMock func(sqlmock.Sqlmock, uuid.UUID)
		expectedError error
		expectedExist bool
	}{
		{
			name:     "success author exists",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `authors` WHERE id = \\?").
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedExist: true,
			expectedError: nil,
		},
		{
			name:     "error author does not exist",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `authors` WHERE id = \\?").
					WithArgs(id).
					WillReturnRows(rows)
			},
			expectedExist: false,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock, test.authorID)

			repo := author.NewAuthorRepository(db, zerolog.Nop())

			exist, err := repo.Exists(context.Background(), test.authorID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedExist, exist)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
