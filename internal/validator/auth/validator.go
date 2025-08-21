package auth

import (
	"fmt"
	"log/slog"
	"sso/protos/gen/go/proto/sso"

	"github.com/go-playground/validator/v10"
)

type LoginDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	AppID    int32  `validate:"required,gt=0"`
}
type RegisterDTO struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}
type IsAdminDTO struct {
	UserID int64 `validate:"required,gt=0"`
}

// Validator — оболочка
type AuthValidator struct {
	log       *slog.Logger
	validator *validator.Validate
}

// New — конструктор
func New(log *slog.Logger, validator *validator.Validate) *AuthValidator {
	return &AuthValidator{
		log:       log,
		validator: validator,
	}
}

// ValidateLogin — возвращает *ValidationError при ошибках валидации
func (v *AuthValidator) ValidateLogin(req *sso.LoginRequest) (validator.ValidationErrors, error) {
	var op = "Validator.ValidateLogin"

	dto := LoginDTO{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		AppID:    req.GetAppId(),
	}
	if err := v.validator.Struct(dto); err != nil {
		return v.logValidationErrors(op, err)
	}
	return nil, nil
}

func (v *AuthValidator) ValidateRegister(req *sso.RegisterRequest) (validator.ValidationErrors, error) {
	var op = "Validator.ValidateRegister"

	dto := RegisterDTO{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := v.validator.Struct(dto); err != nil {
		return v.logValidationErrors(op, err)
	}
	return nil, nil
}

func (v *AuthValidator) ValidateIsAdmin(req *sso.IsAdminRequest) (validator.ValidationErrors, error) {
	var op = "Validator.ValidateIsAdmin"

	dto := IsAdminDTO{
		UserID: req.GetUserId(),
	}

	if err := v.validator.Struct(dto); err != nil {
		return v.logValidationErrors(op, err)
	}
	return nil, nil
}

// logValidationErrors логирует ошибки валидации и возвращает их
func (v *AuthValidator) logValidationErrors(op string, err error) (validator.ValidationErrors, error) {
	log := v.log.With(slog.String("op", op))

	if errs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range errs {
			messages = append(messages,
				fmt.Sprintf("field=%s, rule=%s, param=%s",
					e.Field(), e.Tag(), e.Param()))
		}

		log.Warn("validation failed", slog.Any("errors", messages))
		return errs, nil
	}

	log.Error("unexpected validation error", slog.Any("err", err))
	return nil, err
}
