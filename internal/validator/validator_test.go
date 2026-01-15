package validator_test

import (
	"testing"

	"go-boilerplate-rest-api-chi/internal/validator"
)

func TestNew(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		v := validator.New()
		if v == nil {
			t.Errorf("Expected validator instance, got nil")
		}
	})
}
