package validator

import (
	"fmt"
	"log/slog"
	"sso/internal/validator/auth"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Auth *auth.AuthValidator
}

func New(log *slog.Logger) *Validator {
	validator := validator.New()
	authValidator := auth.New(log, validator)
	return &Validator{
		Auth: authValidator,
	}
}

// HumanizeValidationErrors преобразует ошибки в человекочитаемые сообщения
func (v *Validator) HumanizeValidationErrors(verrs validator.ValidationErrors) []string {
	if verrs == nil {
		return nil
	}

	var messages []string
	for _, fe := range verrs {
		switch fe.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("field '%s' is required", fe.Field()))
		case "email":
			messages = append(messages, fmt.Sprintf("field '%s' must be a valid email", fe.Field()))
		case "min":
			messages = append(messages, fmt.Sprintf("field '%s' must be greater than %s", fe.Field(), fe.Param()))
		default:
			messages = append(messages, fmt.Sprintf("field '%s' is invalid", fe.Field()))
		}
	}

	return messages
}
