package author_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"go-boilerplate-rest-api-chi/internal/author"
	"go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
	"go-boilerplate-rest-api-chi/internal/mocks"
)

func TestAuthorService_CreateAuthor(t *testing.T) {
	tests := []struct {
		name             string
		input            *dto.CreateAuthorRequest
		configureMock    func(*mocks.MockAuthorRepository)
		expectedResponse *entity.Author
		expectedError    error
	}{
		{
			name: "success create author",
			input: &dto.CreateAuthorRequest{
				Name: "J.K. Rowling",
			},
			configureMock: func(mockRepo *mocks.MockAuthorRepository) {
				sampleAuthor := &entity.Author{
					Name: "J.K. Rowling",
				}

				mockRepo.EXPECT().
					Create(gomock.Any(), sampleAuthor).
					Return(&entity.Author{
						ID:   uuid.New(),
						Name: "J.K. Rowling",
					}, nil)
			},
			expectedResponse: &entity.Author{
				Name: "J.K. Rowling",
			},
		},
		{
			name: "error duplicate author",
			input: &dto.CreateAuthorRequest{
				Name: "Duplicate Author",
			},
			configureMock: func(mockRepo *mocks.MockAuthorRepository) {
				expectedEntity := &entity.Author{
					Name: "Duplicate Author",
				}

				mockRepo.EXPECT().
					Create(gomock.Any(), expectedEntity).
					Return(nil, author.ErrDuplicate)
			},
			expectedError: author.ErrDuplicate,
		},
		{
			name: "error database error",
			input: &dto.CreateAuthorRequest{
				Name: "Test Author",
			},
			configureMock: func(mockRepo *mocks.MockAuthorRepository) {
				expectedEntity := &entity.Author{
					Name: "Test Author",
				}

				mockRepo.EXPECT().
					Create(gomock.Any(), expectedEntity).
					Return(nil, errors.New("database connection failed"))
			},
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)

			test.configureMock(authorRepoMock)
			service := author.NewAuthorService(authorRepoMock, zerolog.Nop())

			result, err := service.CreateAuthor(context.Background(), test.input)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResponse.Name, result.Name)
				assert.NotEqual(t, uuid.Nil, result.ID)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestAuthorService_GetAuthorByID(t *testing.T) {
	tests := []struct {
		name             string
		authorID         uuid.UUID
		configureMock    func(*mocks.MockAuthorRepository)
		expectedResponse *entity.Author
		expectedError    error
	}{
		{
			name:     "success create author",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mockRepo *mocks.MockAuthorRepository) {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626")).
					Return(
						&entity.Author{
							ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
							Name: "J.K. Rowling",
						}, nil,
					)
			},
			expectedResponse: &entity.Author{
				ID:   uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
				Name: "J.K. Rowling",
			},
		},
		{
			name:     "error author not found ",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mockRepo *mocks.MockAuthorRepository) {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626")).
					Return(nil, author.ErrNotFound)
			},
			expectedError: author.ErrNotFound,
		},
		{
			name:     "error database error",
			authorID: uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626"),
			configureMock: func(mockRepo *mocks.MockAuthorRepository) {
				mockRepo.EXPECT().
					GetByID(gomock.Any(), uuid.MustParse("eb21d07a-7ab3-40db-bfd3-448093bc5626")).
					Return(nil, errors.New("database connection failed"))
			},
			expectedError: errors.New("database connection failed"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)
			authorRepoMock := mocks.NewMockAuthorRepository(ctrl)

			test.configureMock(authorRepoMock)
			service := author.NewAuthorService(authorRepoMock, zerolog.Nop())

			result, err := service.GetAuthorByID(context.Background(), test.authorID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			if test.expectedResponse != nil {
				assert.NotNil(t, result)
				assert.Equal(t, test.expectedResponse.Name, result.Name)
				assert.NotEqual(t, uuid.Nil, result.ID)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}
