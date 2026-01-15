package book_test

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

	"go-boilerplate-rest-api-chi/internal/book"
	"go-boilerplate-rest-api-chi/internal/entity"
	"go-boilerplate-rest-api-chi/internal/test-utils"
)

func TestBookRepository_Create(t *testing.T) {
	tests := []struct {
		name             string
		input            *entity.Book
		configureMock    func(sqlmock.Sqlmock, *entity.Book)
		expectedError    error
		expectedResponse *entity.Book
	}{
		{
			name: "success create book",
			input: &entity.Book{
				Title:       "Les miserables",
				Description: "Les Misérables raconte la vie de Jean Valjean.",
				AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			},
			configureMock: func(mock sqlmock.Sqlmock, input *entity.Book) {
				mock.ExpectExec("INSERT INTO `books`").
					WithArgs(
						sqlmock.AnyArg(),
						input.Title,
						input.Description,
						input.AuthorID,
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
					).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
			expectedResponse: &entity.Book{
				Title:       "Les miserables",
				Description: "Les Misérables raconte la vie de Jean Valjean.",
				AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			},
		},
		{
			name: "error duplicate book",
			input: &entity.Book{
				Title:       "Duplicate Book",
				Description: "Duplicate book description",
				AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			},
			configureMock: func(mock sqlmock.Sqlmock, input *entity.Book) {
				mock.ExpectExec("INSERT INTO `books`").
					WithArgs(
						sqlmock.AnyArg(), // ID
						input.Title,
						input.Description,
						input.AuthorID,
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
					).WillReturnError(gorm.ErrDuplicatedKey)
			},
			expectedError:    book.ErrDuplicate,
			expectedResponse: nil,
		},
		{
			name: "error database connection failed",
			input: &entity.Book{
				Title:       "Les miserables",
				Description: "Les Misérables raconte la vie de Jean Valjean.",
				AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			},
			configureMock: func(mock sqlmock.Sqlmock, input *entity.Book) {
				mock.ExpectExec("INSERT INTO `books`").
					WithArgs(
						sqlmock.AnyArg(), // ID
						input.Title,
						input.Description,
						input.AuthorID,
						sqlmock.AnyArg(), // CreatedAt
						sqlmock.AnyArg(), // UpdatedAt
					).WillReturnError(gorm.ErrInvalidDB)
			},
			expectedError:    gorm.ErrInvalidDB,
			expectedResponse: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock, test.input)

			repo := book.NewBookRepository(db, zerolog.Nop())

			newBook, err := repo.Create(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, newBook)
				assert.NotEqual(t, uuid.Nil, newBook.ID)
				assert.Equal(t, test.expectedResponse.Title, newBook.Title)
				assert.Equal(t, test.expectedResponse.Description, newBook.Description)
				assert.Equal(t, test.expectedResponse.AuthorID, newBook.AuthorID)
			} else {
				assert.Nil(t, newBook)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookRepository_GetAll(t *testing.T) {
	tests := []struct {
		name             string
		configureMock    func(sqlmock.Sqlmock)
		expectedError    error
		expectedResponse []*entity.Book
	}{
		{
			name: "success get all books",
			configureMock: func(mock sqlmock.Sqlmock) {
				now := time.Now()
				authorID := uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626")
				bookID := uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0")
				bookID2 := uuid.MustParse("b1c2d3e4-f5a6-7890-1234-56789abcdef1")
				bookID3 := uuid.MustParse("c1d2e3f4-a5b6-7890-1234-56789abcdef2")

				rows := sqlmock.NewRows([]string{"id", "title", "description", "author_id", "created_at", "updated_at"}).
					AddRow(bookID, "Book One", "Description One", authorID, now, now).
					AddRow(bookID2, "Book Two", "Description Two", authorID, now, now).
					AddRow(bookID3, "Book Three", "Description Three", authorID, now, now)

				mock.ExpectQuery("SELECT \\* FROM `books`").
					WillReturnRows(rows)

				authorRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(authorID, "Victor Hugo", now, now)

				mock.ExpectQuery("SELECT \\* FROM `authors` WHERE `authors`.`id` = \\?").
					WithArgs(authorID).
					WillReturnRows(authorRows)
			},
			expectedError: nil,
			expectedResponse: []*entity.Book{
				{
					ID:          uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
					Title:       "Book One",
					Description: "Description One",
					AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
					Author: &entity.Author{
						ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
						Name: "Victor Hugo",
					},
				},
				{
					ID:          uuid.MustParse("b1c2d3e4-f5a6-7890-1234-56789abcdef1"),
					Title:       "Book Two",
					Description: "Description Two",
					AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
					Author: &entity.Author{
						ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
						Name: "Victor Hugo",
					},
				},
				{
					ID:          uuid.MustParse("c1d2e3f4-a5b6-7890-1234-56789abcdef2"),
					Title:       "Book Three",
					Description: "Description Three",
					AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
					Author: &entity.Author{
						ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
						Name: "Victor Hugo",
					},
				},
			},
		},
		{
			name: "error database connection failed",
			configureMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT \\* FROM `books`").
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectedError:    gorm.ErrInvalidDB,
			expectedResponse: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock)

			repo := book.NewBookRepository(db, zerolog.Nop())

			books, err := repo.GetAll(context.Background())

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, books)
				assert.Len(t, books, len(test.expectedResponse))
				for i := range books {
					assert.Equal(t, test.expectedResponse[i].ID, books[i].ID)
					assert.Equal(t, test.expectedResponse[i].Title, books[i].Title)
					assert.Equal(t, test.expectedResponse[i].Description, books[i].Description)
					assert.Equal(t, test.expectedResponse[i].AuthorID, books[i].AuthorID)
					if test.expectedResponse[i].Author != nil {
						assert.NotNil(t, books[i].Author)
						assert.Equal(t, test.expectedResponse[i].Author.ID, books[i].Author.ID)
						assert.Equal(t, test.expectedResponse[i].Author.Name, books[i].Author.Name)
					} else {
						assert.Nil(t, books[i].Author)
					}
				}
			} else {
				assert.Nil(t, books)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookRepository_GetByID(t *testing.T) {
	tests := []struct {
		name             string
		bookID           uuid.UUID
		configureMock    func(sqlmock.Sqlmock, uuid.UUID)
		expectedError    error
		expectedResponse *entity.Book
	}{
		{
			name:   "success get book by id",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				now := time.Now()
				authorID := uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626")

				booksRows := sqlmock.NewRows([]string{"id", "title", "description", "author_id", "created_at", "updated_at"}).
					AddRow(id, "Book One", "Description One", authorID, now, now)

				mock.ExpectQuery("SELECT \\* FROM `books` WHERE id = \\? ORDER BY `books`.`id` LIMIT \\?").
					WithArgs(id, 1).
					WillReturnRows(booksRows)

				authorRows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
					AddRow(authorID, "Victor Hugo", now, now)

				mock.ExpectQuery("SELECT \\* FROM `authors` WHERE `authors`.`id` = \\?").
					WithArgs(authorID).
					WillReturnRows(authorRows)

			},
			expectedError: nil,
			expectedResponse: &entity.Book{
				ID:          uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
				Title:       "Book One",
				Description: "Description One",
				AuthorID:    uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
				Author: &entity.Author{
					ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
					Name: "Victor Hugo",
				},
			},
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery("SELECT \\* FROM `books` WHERE id = \\? ORDER BY `books`.`id` LIMIT \\?").
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)

			},
			expectedError:    book.ErrNotFound,
			expectedResponse: nil,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			configureMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery("SELECT \\* FROM `books` WHERE id = \\? ORDER BY `books`.`id` LIMIT \\?").
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
			test.configureMock(mock, test.bookID)

			repo := book.NewBookRepository(db, zerolog.Nop())

			book, err := repo.GetByID(context.Background(), test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, book)
				assert.Equal(t, test.expectedResponse.ID, book.ID)
				assert.Equal(t, test.expectedResponse.Title, book.Title)
				assert.Equal(t, test.expectedResponse.Description, book.Description)
				assert.Equal(t, test.expectedResponse.AuthorID, book.AuthorID)
				assert.Equal(t, test.expectedResponse.Author.ID, book.Author.ID)
				assert.Equal(t, test.expectedResponse.Author.Name, book.Author.Name)
			} else {
				assert.Nil(t, book)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookRepository_Update(t *testing.T) {
	tests := []struct {
		name          string
		bookID        uuid.UUID
		updates       map[string]interface{}
		configureMock func(sqlmock.Sqlmock, uuid.UUID, map[string]interface{})
		expectedError error
	}{
		{
			name:   "success update book",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			updates: map[string]interface{}{
				"description": "Updated description",
			},
			configureMock: func(mock sqlmock.Sqlmock, bookID uuid.UUID, updates map[string]interface{}) {
				mock.ExpectExec("UPDATE `books` SET `description`=\\?,`updated_at`=\\? WHERE `id` = \\?").
					WithArgs(
						updates["description"],
						sqlmock.AnyArg(), // updated_at
						bookID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: nil,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			updates: map[string]interface{}{
				"description": "Updated description",
			},
			configureMock: func(mock sqlmock.Sqlmock, bookID uuid.UUID, updates map[string]interface{}) {
				mock.ExpectExec("UPDATE `books` SET `description`=\\?,`updated_at`=\\? WHERE `id` = \\?").
					WithArgs(
						updates["description"],
						sqlmock.AnyArg(),
						bookID,
					).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectedError: gorm.ErrInvalidDB,
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			updates: map[string]interface{}{
				"description": "Updated description",
			},
			configureMock: func(mock sqlmock.Sqlmock, bookID uuid.UUID, updates map[string]interface{}) {
				mock.ExpectExec("UPDATE `books` SET `description`=\\?,`updated_at`=\\? WHERE `id` = \\?").
					WithArgs(
						updates["description"],
						sqlmock.AnyArg(),
						bookID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: book.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock, test.bookID, test.updates)

			repo := book.NewBookRepository(db, zerolog.Nop())

			err := repo.Update(context.Background(), test.bookID, test.updates)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBookRepository_Delete(t *testing.T) {
	tests := []struct {
		name          string
		bookID        uuid.UUID
		configureMock func(sqlmock.Sqlmock, uuid.UUID)
		expectedError error
	}{
		{
			name:   "success delete book",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			configureMock: func(mock sqlmock.Sqlmock, bookID uuid.UUID) {
				mock.ExpectExec("DELETE FROM `books` WHERE `books`.`id` = \\?").
					WithArgs(bookID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			configureMock: func(mock sqlmock.Sqlmock, bookID uuid.UUID) {
				mock.ExpectExec("DELETE FROM `books` WHERE `books`.`id` = \\?").
					WithArgs(bookID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: book.ErrNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("a1b2c3d4-e5f6-7890-1234-56789abcdef0"),
			configureMock: func(mock sqlmock.Sqlmock, bookID uuid.UUID) {
				mock.ExpectExec("DELETE FROM `books` WHERE `books`.`id` = \\?").
					WithArgs(bookID).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectedError: gorm.ErrInvalidDB,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := testutils.NewGormMySQL(t)
			test.configureMock(mock, test.bookID)

			repo := book.NewBookRepository(db, zerolog.Nop())

			err := repo.Delete(context.Background(), test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
