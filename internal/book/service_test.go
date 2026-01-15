package book_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/book"
	"go-boilerplate-rest-api-chi/internal/book/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
	"go-boilerplate-rest-api-chi/internal/mocks"
)

func TestBookService_CreateBook(t *testing.T) {
	tests := []struct {
		name             string
		input            *dto.CreateBookRequest
		configureMock    func(*mocks.MockBookRepository, *mocks.MockAuthorRepository)
		expectedResponse *entity.Book
		expectedError    error
	}{
		{
			name: "success create book",
			input: &dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockBookRepository *mocks.MockBookRepository, mockAuthorRepository *mocks.MockAuthorRepository) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockAuthorRepository.EXPECT().
					Exists(gomock.Any(), authorID).
					Return(true, nil)

				sampleBook := &entity.Book{
					Title:       "Book1",
					Description: "Description of book1",
					AuthorID:    authorID,
				}

				mockBookRepository.EXPECT().
					Create(gomock.Any(), sampleBook).
					Return(&entity.Book{
						ID:          uuid.New(),
						Title:       "Book1",
						Description: "Description of book1",
						AuthorID:    authorID,
					}, nil)
			},
			expectedResponse: &entity.Book{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4"),
			},
			expectedError: nil,
		},
		{
			name: "error invalid AuthorID",
			input: &dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "invalid uuid",
			},
			configureMock:    func(mockBookRepository *mocks.MockBookRepository, mockAuthorRepository *mocks.MockAuthorRepository) {},
			expectedResponse: nil,
			expectedError:    book.ErrInvalidAuthorId,
		},
		{
			name: "error exists check fails",
			input: &dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockBookRepository *mocks.MockBookRepository, mockAuthorRepository *mocks.MockAuthorRepository) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockAuthorRepository.EXPECT().
					Exists(gomock.Any(), authorID).
					Return(false, gorm.ErrInvalidDB)
			},
			expectedResponse: nil,
			expectedError:    gorm.ErrInvalidDB,
		},
		{
			name: "error author not found",
			input: &dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockBookRepository *mocks.MockBookRepository, mockAuthorRepository *mocks.MockAuthorRepository) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockAuthorRepository.EXPECT().
					Exists(gomock.Any(), authorID).
					Return(false, nil)
			},
			expectedResponse: nil,
			expectedError:    author.ErrNotFound,
		},
		{
			name: "error duplicate book",
			input: &dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockBookRepository *mocks.MockBookRepository, mockAuthorRepository *mocks.MockAuthorRepository) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockAuthorRepository.EXPECT().
					Exists(gomock.Any(), authorID).
					Return(true, nil)

				sampleBook := &entity.Book{
					Title:       "Book1",
					Description: "Description of book1",
					AuthorID:    authorID,
				}

				mockBookRepository.EXPECT().
					Create(gomock.Any(), sampleBook).
					Return(nil, book.ErrDuplicate)
			},
			expectedResponse: nil,
			expectedError:    book.ErrDuplicate,
		},
		{
			name: "error database connection failed",
			input: &dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description of book1",
				AuthorID:    "779404e4-2660-4c80-b958-cfa72515e7d4",
			},
			configureMock: func(mockBookRepository *mocks.MockBookRepository, mockAuthorRepository *mocks.MockAuthorRepository) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				mockAuthorRepository.EXPECT().
					Exists(gomock.Any(), authorID).
					Return(false, gorm.ErrInvalidDB)
			},
			expectedResponse: nil,
			expectedError:    gorm.ErrInvalidDB,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			test.configureMock(bookRepoMock, authorRepoMock)
			service := book.NewBookService(bookRepoMock, authorRepoMock, zerolog.Nop())

			result, err := service.CreateBook(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, result)
				assert.NotEqual(t, uuid.Nil, result.ID)
				assert.Equal(t, test.expectedResponse.Title, result.Title)
				assert.Equal(t, test.expectedResponse.Description, result.Description)
				assert.Equal(t, test.expectedResponse.AuthorID, result.AuthorID)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestBookService_GetAllBooks(t *testing.T) {
	tests := []struct {
		name             string
		configureMock    func(*mocks.MockBookRepository)
		expectedResponse []*entity.Book
		expectedError    error
	}{
		{
			name: "success get all books",
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				authorID := uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4")

				bookRepository.EXPECT().
					GetAll(gomock.Any()).
					Return([]*entity.Book{
						{
							ID:          uuid.MustParse("619c69fb-9bcb-451e-b825-29b81697a531"),
							Title:       "Book1",
							Description: "Description of book1",
							Author: &entity.Author{
								ID:   authorID,
								Name: "Author1",
							},
						},
						{
							ID:          uuid.MustParse("92e37eac-0097-405e-a5cc-7055a39bdde8"),
							Title:       "Book2",
							Description: "Description of book2",
							Author: &entity.Author{
								ID:   authorID,
								Name: "Author1",
							},
						},
						{
							ID:          uuid.MustParse("dec6b550-292e-45ad-b269-f59bfa06bf01"),
							Title:       "Book3",
							Description: "Description of book3",
							Author: &entity.Author{
								ID:   authorID,
								Name: "Author1",
							},
						},
					}, nil)
			},
			expectedResponse: []*entity.Book{
				{
					ID:          uuid.MustParse("619c69fb-9bcb-451e-b825-29b81697a531"),
					Title:       "Book1",
					Description: "Description of book1",
					Author: &entity.Author{
						ID:   uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4"),
						Name: "Author1",
					},
				},
				{
					ID:          uuid.MustParse("92e37eac-0097-405e-a5cc-7055a39bdde8"),
					Title:       "Book2",
					Description: "Description of book2",
					Author: &entity.Author{
						ID:   uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4"),
						Name: "Author1",
					},
				},
				{
					ID:          uuid.MustParse("dec6b550-292e-45ad-b269-f59bfa06bf01"),
					Title:       "Book3",
					Description: "Description of book3",
					Author: &entity.Author{
						ID:   uuid.MustParse("779404e4-2660-4c80-b958-cfa72515e7d4"),
						Name: "Author1",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "error no book found",
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookRepository.EXPECT().
					GetAll(gomock.Any()).
					Return([]*entity.Book{}, nil)
			},
			expectedResponse: nil,
			expectedError:    book.ErrNotFound,
		},
		{
			name: "error database connection failed",
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookRepository.EXPECT().
					GetAll(gomock.Any()).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectedResponse: nil,
			expectedError:    gorm.ErrInvalidDB,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			test.configureMock(bookRepoMock)
			service := book.NewBookService(bookRepoMock, authorRepoMock, zerolog.Nop())

			books, err := service.GetAllBooks(context.Background())

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, books)
				for i := range books {
					assert.Equal(t, test.expectedResponse[i].ID, books[i].ID)
					assert.Equal(t, test.expectedResponse[i].Title, books[i].Title)
					assert.Equal(t, test.expectedResponse[i].Description, books[i].Description)
					assert.Equal(t, test.expectedResponse[i].Author.ID, books[i].Author.ID)
					assert.Equal(t, test.expectedResponse[i].Author.Name, books[i].Author.Name)
				}
			} else {
				assert.Nil(t, books)
			}
		})
	}
}

func TestBookService_GetBookByID(t *testing.T) {
	tests := []struct {
		name             string
		bookID           uuid.UUID
		configureMock    func(*mocks.MockBookRepository)
		expectedResponse *entity.Book
		expectedError    error
	}{
		{
			name:   "success get book by ID",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookRepository.EXPECT().
					GetByID(gomock.Any(), uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")).
					Return(&entity.Book{
						ID:          uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
						Title:       "Book1",
						Description: "Description of book1",
						Author: &entity.Author{
							ID:   uuid.MustParse("b57063a5-409c-457c-bbec-d5850b2e3761"),
							Name: "Author1",
						},
					}, nil)
			},
			expectedResponse: &entity.Book{
				ID:          uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
				Title:       "Book1",
				Description: "Description of book1",
				Author: &entity.Author{
					ID:   uuid.MustParse("b57063a5-409c-457c-bbec-d5850b2e3761"),
					Name: "Author1",
				},
			},
			expectedError: nil,
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookRepository.EXPECT().
					GetByID(gomock.Any(), uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")).
					Return(nil, book.ErrNotFound)
			},
			expectedResponse: nil,
			expectedError:    book.ErrNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookRepository.EXPECT().
					GetByID(gomock.Any(), uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")).
					Return(nil, gorm.ErrInvalidDB)
			},
			expectedResponse: nil,
			expectedError:    gorm.ErrInvalidDB,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			test.configureMock(bookRepoMock)
			service := book.NewBookService(bookRepoMock, authorRepoMock, zerolog.Nop())

			books, err := service.GetBookByID(context.Background(), test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, books)
				assert.NotEqual(t, test.expectedResponse.ID, uuid.Nil)
				assert.Equal(t, test.expectedResponse.Title, books.Title)
				assert.Equal(t, test.expectedResponse.Description, books.Description)
				assert.Equal(t, test.expectedResponse.Author.ID, books.Author.ID)
				assert.Equal(t, test.expectedResponse.Author.Name, books.Author.Name)
			} else {
				assert.Nil(t, books)
			}
		})
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	tests := []struct {
		name          string
		bookID        uuid.UUID
		input         *dto.UpdateBookRequest
		configureMock func(*mocks.MockBookRepository)
		expectedError error
	}{
		{
			name:   "success update book",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			input: &dto.UpdateBookRequest{
				Description: "New updated description",
			},
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookID := uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")

				updates := map[string]interface{}{
					"description": "New updated description",
				}

				bookRepository.EXPECT().
					Update(gomock.Any(), bookID, updates).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			input: &dto.UpdateBookRequest{
				Description: "Updated description",
			},
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookID := uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")

				bookRepository.EXPECT().
					Update(gomock.Any(), bookID, gomock.Any()).
					Return(book.ErrNotFound)
			},
			expectedError: book.ErrNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			input: &dto.UpdateBookRequest{
				Description: "Updated description",
			},
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookID := uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")

				bookRepository.EXPECT().
					Update(gomock.Any(), bookID, gomock.Any()).
					Return(gorm.ErrInvalidDB)
			},
			expectedError: gorm.ErrInvalidDB,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			test.configureMock(bookRepoMock)
			service := book.NewBookService(bookRepoMock, authorRepoMock, zerolog.Nop())

			err := service.UpdateBook(context.Background(), test.input, test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookService_DeleteBook(t *testing.T) {
	tests := []struct {
		name          string
		bookID        uuid.UUID
		configureMock func(*mocks.MockBookRepository)
		expectedError error
	}{
		{
			name:   "success delete book",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookID := uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")

				bookRepository.EXPECT().
					Delete(gomock.Any(), bookID).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error book not found",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookID := uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")

				bookRepository.EXPECT().
					Delete(gomock.Any(), bookID).
					Return(book.ErrNotFound)
			},
			expectedError: book.ErrNotFound,
		},
		{
			name:   "error database connection failed",
			bookID: uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06"),
			configureMock: func(bookRepository *mocks.MockBookRepository) {
				bookID := uuid.MustParse("c6efb683-455d-4ca0-b8aa-83ca9b930a06")

				bookRepository.EXPECT().
					Delete(gomock.Any(), bookID).
					Return(gorm.ErrInvalidDB)
			},
			expectedError: gorm.ErrInvalidDB,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)
			bookRepoMock := mocks.NewMockBookRepository(ctrl)

			test.configureMock(bookRepoMock)
			service := book.NewBookService(bookRepoMock, authorRepoMock, zerolog.Nop())

			err := service.DeleteBook(context.Background(), test.bookID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
