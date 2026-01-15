package book_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"go-boilerplate-rest-api-chi/internal/author"
	authorDTO "go-boilerplate-rest-api-chi/internal/author/dto"
	"go-boilerplate-rest-api-chi/internal/book"
	"go-boilerplate-rest-api-chi/internal/book/dto"
	"go-boilerplate-rest-api-chi/internal/entity"
	"go-boilerplate-rest-api-chi/internal/mocks"
	"go-boilerplate-rest-api-chi/internal/response"
	"go-boilerplate-rest-api-chi/internal/validator"
)

func TestBookHandler_CreateBook(t *testing.T) {
	tests := []struct {
		name               string
		requestBody        interface{}
		configureMock      func(service *mocks.MockBookService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "success create book",
			requestBody: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.CreateBookRequest{
					Title:       "Book1",
					Description: "Description1",
					AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
				}

				mockService.EXPECT().
					CreateBook(gomock.Any(), input).
					Return(&entity.Book{
						ID:          uuid.MustParse("13867a7d-d1c4-4a06-aa60-42741a4fbbbd"),
						Title:       "Book1",
						Description: "Description1",
						Author: &entity.Author{
							ID:   uuid.MustParse("24319e61-32d0-49f3-987f-019b734ed9c7"),
							Name: "Author1",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse: &book.BookSuccessResponse{
				Status:  "success",
				Message: "Book created successfully",
				Book: &dto.BookResponse{
					ID:          "13867a7d-d1c4-4a06-aa60-42741a4fbbbd",
					Title:       "Book1",
					Description: "Description1",
					Author: authorDTO.AuthorResponse{
						ID:   "24319e61-32d0-49f3-987f-019b734ed9c7",
						Name: "Author1",
					},
				},
			},
		},
		{
			name:               "error invalid json",
			requestBody:        nil,
			configureMock:      func(service *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Invalid request body",
			},
		},
		{
			name: "error validation fails empty name",
			requestBody: dto.CreateBookRequest{
				Title:       "",
				Description: "",
				AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
			},
			configureMock:      func(service *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: response.ValidationErrorResponse{
				Status:  "error",
				Message: "Validation failed",
				Errors: []response.ValidationErrorDetail{{
					Field:   "Title",
					Message: "Title is required",
				},
					{
						Field:   "Description",
						Message: "Description is required",
					}},
			},
		},
		{
			name: "error service internal error",
			requestBody: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.CreateBookRequest{
					Title:       "Book1",
					Description: "Description1",
					AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
				}
				mockService.EXPECT().
					CreateBook(gomock.Any(), input).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
			},
		},
		{
			name: "error invalid author id",
			requestBody: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    "invalid-uuid",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.CreateBookRequest{
					Title:       "Book1",
					Description: "Description1",
					AuthorID:    "invalid-uuid",
				}
				mockService.EXPECT().
					CreateBook(gomock.Any(), input).
					Return(nil, book.ErrInvalidAuthorId)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "invalid author ID",
			},
		},
		{
			name: "error author not found",
			requestBody: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.CreateBookRequest{
					Title:       "Book1",
					Description: "Description1",
					AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
				}
				mockService.EXPECT().
					CreateBook(gomock.Any(), input).
					Return(nil, author.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Author not found",
			},
		},
		{
			name: "error duplicate book",
			requestBody: dto.CreateBookRequest{
				Title:       "Book1",
				Description: "Description1",
				AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.CreateBookRequest{
					Title:       "Book1",
					Description: "Description1",
					AuthorID:    "24319e61-32d0-49f3-987f-019b734ed9c7",
				}
				mockService.EXPECT().
					CreateBook(gomock.Any(), input).
					Return(nil, book.ErrDuplicate)
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Book with this name already exists",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockBookService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := book.NewBookHandler(mockService, v, zerolog.Nop())

			var body *bytes.Buffer
			if test.requestBody == nil {
				body = bytes.NewBuffer([]byte{})
			} else {
				b, err := json.Marshal(test.requestBody)
				require.NoError(t, err)
				body = bytes.NewBuffer(b)
			}

			req := httptest.NewRequest(http.MethodPost, "/books", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/books", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())

		})
	}
}

func TestBookHandler_GetAllBooks(t *testing.T) {
	tests := []struct {
		name               string
		configureMock      func(service *mocks.MockBookService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "success get all books",
			configureMock: func(mockService *mocks.MockBookService) {
				author := entity.Author{
					ID:   uuid.MustParse("24319e61-32d0-49f3-987f-019b734ed9c7"),
					Name: "Author1",
				}
				mockService.EXPECT().
					GetAllBooks(gomock.Any()).
					Return([]*entity.Book{
						{
							ID:          uuid.MustParse("13867a7d-d1c4-4a06-aa60-42741a4fbbbd"),
							Title:       "Book1",
							Description: "Description1",
							Author:      &author,
						},
						{
							ID:          uuid.MustParse("66509608-3ca2-46d0-99d6-8ad989fe0061"),
							Title:       "Book2",
							Description: "Description2",
							Author:      &author,
						},
						{
							ID:          uuid.MustParse("d58905b0-1d21-47ee-805f-ccc92aba2453"),
							Title:       "Book3",
							Description: "Description3",
							Author:      &author,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: book.BooksSuccessResponse{
				Status:  "success",
				Message: "Books retrieved successfully",
				Books: []dto.BookResponse{
					{
						ID:          "13867a7d-d1c4-4a06-aa60-42741a4fbbbd",
						Title:       "Book1",
						Description: "Description1",
						Author: authorDTO.AuthorResponse{
							ID:   "24319e61-32d0-49f3-987f-019b734ed9c7",
							Name: "Author1",
						},
					},
					{
						ID:          "66509608-3ca2-46d0-99d6-8ad989fe0061",
						Title:       "Book2",
						Description: "Description2",
						Author: authorDTO.AuthorResponse{
							ID:   "24319e61-32d0-49f3-987f-019b734ed9c7",
							Name: "Author1",
						},
					},
					{
						ID:          "d58905b0-1d21-47ee-805f-ccc92aba2453",
						Title:       "Book3",
						Description: "Description3",
						Author: authorDTO.AuthorResponse{
							ID:   "24319e61-32d0-49f3-987f-019b734ed9c7",
							Name: "Author1",
						},
					},
				},
			},
		},
		{
			name: "error service internal error",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					GetAllBooks(gomock.Any()).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockBookService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := book.NewBookHandler(mockService, v, zerolog.Nop())

			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/books", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())

		})
	}
}

func TestBookHandler_GetBookByID(t *testing.T) {
	tests := []struct {
		name               string
		idUrlParam         string
		configureMock      func(service *mocks.MockBookService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:       "success get book by id",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					GetBookByID(gomock.Any(), uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(&entity.Book{
						ID:          uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227"),
						Title:       "Book1",
						Description: "Description1",
						Author: &entity.Author{
							ID:   uuid.MustParse("88a49625-ee9d-456d-9541-e359454eb40c"),
							Name: "Author1",
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &book.BookSuccessResponse{
				Status:  "success",
				Message: "Book retrieved successfully",
				Book: &dto.BookResponse{
					ID:          "3a310074-b63f-455e-996f-63a5afffc227",
					Title:       "Book1",
					Description: "Description1",
					Author: authorDTO.AuthorResponse{
						ID:   "88a49625-ee9d-456d-9541-e359454eb40c",
						Name: "Author1",
					},
				},
			},
		},
		{
			name:       "error book not found",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					GetBookByID(gomock.Any(), uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(nil, book.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Book not found",
			},
		},
		{
			name:               "error invalid uuid",
			idUrlParam:         "invalid-uuid",
			configureMock:      func(mockService *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Invalid uuid",
			},
		},
		{
			name:       "error service internal error",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					GetBookByID(gomock.Any(), uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(nil, errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockBookService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := book.NewBookHandler(mockService, v, zerolog.Nop())

			url := fmt.Sprintf("/books/%s", test.idUrlParam)
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/books", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())

		})
	}
}

func TestBookHandler_UpdateBook(t *testing.T) {
	tests := []struct {
		name               string
		idUrlParam         string
		requestBody        interface{}
		configureMock      func(service *mocks.MockBookService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:       "success update book",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			requestBody: dto.UpdateBookRequest{
				Description: "Updated Description",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.UpdateBookRequest{
					Description: "Updated Description",
				}

				mockService.EXPECT().
					UpdateBook(gomock.Any(), input, uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &response.SuccessResponse{
				Status:  "success",
				Message: "Book updated successfully",
			},
		},
		{
			name:       "error invalid uuid",
			idUrlParam: "invalid-uuid",
			requestBody: dto.UpdateBookRequest{
				Description: "Updated Description",
			},
			configureMock:      func(mockService *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Invalid uuid",
			},
		},
		{
			name:               "error invalid json",
			idUrlParam:         "3a310074-b63f-455e-996f-63a5afffc227",
			requestBody:        nil,
			configureMock:      func(mockService *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Invalid request body",
			},
		},
		{
			name:       "error validation fails empty description",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			requestBody: dto.UpdateBookRequest{
				Description: "",
			},
			configureMock:      func(mockService *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: response.ValidationErrorResponse{
				Status:  "error",
				Message: "Validation failed",
				Errors: []response.ValidationErrorDetail{{
					Field:   "Description",
					Message: "Description is required",
				}},
			},
		},
		{
			name:       "error service internal error",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			requestBody: dto.UpdateBookRequest{
				Description: "Updated Description",
			},
			configureMock: func(mockService *mocks.MockBookService) {
				input := &dto.UpdateBookRequest{
					Description: "Updated Description",
				}
				mockService.EXPECT().
					UpdateBook(gomock.Any(), input, uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockBookService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := book.NewBookHandler(mockService, v, zerolog.Nop())

			var body *bytes.Buffer
			if test.requestBody == nil {
				body = bytes.NewBuffer([]byte{})
			} else {
				b, err := json.Marshal(test.requestBody)
				require.NoError(t, err)
				body = bytes.NewBuffer(b)
			}

			url := fmt.Sprintf("/books/%s", test.idUrlParam)
			req := httptest.NewRequest(http.MethodPatch, url, body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/books", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())

		})
	}
}

func TestBookHandler_DeleteBook(t *testing.T) {
	tests := []struct {
		name               string
		idUrlParam         string
		configureMock      func(service *mocks.MockBookService)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:       "success delete book",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					DeleteBook(gomock.Any(), uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: &response.SuccessResponse{
				Status:  "success",
				Message: "Book deleted successfully",
			},
		},
		{
			name:               "error invalid uuid",
			idUrlParam:         "invalid-uuid",
			configureMock:      func(mockService *mocks.MockBookService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Invalid uuid",
			},
		},
		{
			name:       "error book not found",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					DeleteBook(gomock.Any(), uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(book.ErrNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Book not found",
			},
		},
		{
			name:       "error service internal error",
			idUrlParam: "3a310074-b63f-455e-996f-63a5afffc227",
			configureMock: func(mockService *mocks.MockBookService) {
				mockService.EXPECT().
					DeleteBook(gomock.Any(), uuid.MustParse("3a310074-b63f-455e-996f-63a5afffc227")).
					Return(errors.New("database connection failed"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: &response.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mockService := mocks.NewMockBookService(ctrl)
			test.configureMock(mockService)

			v := validator.New()
			handler := book.NewBookHandler(mockService, v, zerolog.Nop())

			url := fmt.Sprintf("/books/%s", test.idUrlParam)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Mount("/books", handler.Routes())

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			expectedJSON, err := json.Marshal(test.expectedResponse)
			require.NoError(t, err)

			assert.JSONEq(t, string(expectedJSON), w.Body.String())

		})
	}
}
