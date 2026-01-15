package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
}

type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"An error occurred"`
}

type ValidationErrorDetail struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email is required"`
}

type ValidationErrorResponse struct {
	Status  string                  `json:"status" example:"error"`
	Message string                  `json:"message" example:"Validation failed"`
	Errors  []ValidationErrorDetail `json:"errors"`
}

func JSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func Success(w http.ResponseWriter, message string) {
	JSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: message,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func ValidationError(w http.ResponseWriter, errors []ValidationErrorDetail) {
	JSON(w, http.StatusBadRequest, ValidationErrorResponse{
		Status:  "error",
		Message: "Validation failed",
		Errors:  errors,
	})
}
