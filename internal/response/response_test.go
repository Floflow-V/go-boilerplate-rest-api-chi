package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go-boilerplate-rest-api-chi/internal/response"
)

func TestJSON(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		data           any
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "success with struct",
			statusCode:     http.StatusOK,
			data:           map[string]string{"message": "test"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"test"}`,
		},
		{
			name:           "success with status created",
			statusCode:     http.StatusCreated,
			data:           map[string]string{"id": "123"},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"123"}`,
		},
		{
			name:           "success with complex data",
			statusCode:     http.StatusOK,
			data:           map[string]interface{}{"count": 5, "active": true},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"active":true,"count":5}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			response.JSON(w, test.statusCode, test.data)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.JSONEq(t, test.expectedBody, w.Body.String())
		})
	}
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		expectedBody response.SuccessResponse
	}{
		{
			name:    "success with custom message",
			message: "Operation completed successfully",
			expectedBody: response.SuccessResponse{
				Status:  "success",
				Message: "Operation completed successfully",
			},
		},
		{
			name:    "success with different message",
			message: "Data created",
			expectedBody: response.SuccessResponse{
				Status:  "success",
				Message: "Data created",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			response.Success(w, test.message)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var result response.SuccessResponse
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)

			assert.Equal(t, test.expectedBody.Status, result.Status)
			assert.Equal(t, test.expectedBody.Message, result.Message)
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		message      string
		expectedBody response.ErrorResponse
	}{
		{
			name:    "error bad request",
			status:  http.StatusBadRequest,
			message: "Invalid input",
			expectedBody: response.ErrorResponse{
				Status:  "error",
				Message: "Invalid input",
			},
		},
		{
			name:    "error not found",
			status:  http.StatusNotFound,
			message: "Resource not found",
			expectedBody: response.ErrorResponse{
				Status:  "error",
				Message: "Resource not found",
			},
		},
		{
			name:    "error internal server",
			status:  http.StatusInternalServerError,
			message: "Internal server error",
			expectedBody: response.ErrorResponse{
				Status:  "error",
				Message: "Internal server error",
			},
		},
		{
			name:    "error conflict",
			status:  http.StatusConflict,
			message: "Resource already exists",
			expectedBody: response.ErrorResponse{
				Status:  "error",
				Message: "Resource already exists",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			response.Error(w, test.status, test.message)

			assert.Equal(t, test.status, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var result response.ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)

			assert.Equal(t, test.expectedBody.Status, result.Status)
			assert.Equal(t, test.expectedBody.Message, result.Message)
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name         string
		errors       []response.ValidationErrorDetail
		expectedBody response.ValidationErrorResponse
	}{
		{
			name: "single validation error",
			errors: []response.ValidationErrorDetail{
				{
					Field:   "email",
					Message: "Email is required",
				},
			},
			expectedBody: response.ValidationErrorResponse{
				Status:  "error",
				Message: "Validation failed",
				Errors: []response.ValidationErrorDetail{
					{
						Field:   "email",
						Message: "Email is required",
					},
				},
			},
		},
		{
			name: "multiple validation errors",
			errors: []response.ValidationErrorDetail{
				{
					Field:   "email",
					Message: "Email is required",
				},
				{
					Field:   "password",
					Message: "Password must be at least 8 characters",
				},
				{
					Field:   "username",
					Message: "Username is already taken",
				},
			},
			expectedBody: response.ValidationErrorResponse{
				Status:  "error",
				Message: "Validation failed",
				Errors: []response.ValidationErrorDetail{
					{
						Field:   "email",
						Message: "Email is required",
					},
					{
						Field:   "password",
						Message: "Password must be at least 8 characters",
					},
					{
						Field:   "username",
						Message: "Username is already taken",
					},
				},
			},
		},
		{
			name:   "empty validation errors",
			errors: []response.ValidationErrorDetail{},
			expectedBody: response.ValidationErrorResponse{
				Status:  "error",
				Message: "Validation failed",
				Errors:  []response.ValidationErrorDetail{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			response.ValidationError(w, test.errors)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var result response.ValidationErrorResponse
			err := json.NewDecoder(w.Body).Decode(&result)
			require.NoError(t, err)

			assert.Equal(t, test.expectedBody.Status, result.Status)
			assert.Equal(t, test.expectedBody.Message, result.Message)
			assert.Len(t, result.Errors, len(test.expectedBody.Errors))

			for i := range result.Errors {
				assert.Equal(t, test.expectedBody.Errors[i].Field, result.Errors[i].Field)
				assert.Equal(t, test.expectedBody.Errors[i].Message, result.Errors[i].Message)
			}
		})
	}
}
