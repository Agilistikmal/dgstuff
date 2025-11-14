package app

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewUnauthorizedError() *AppError {
	return NewAppError(http.StatusUnauthorized, "unauthorized")
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}

func NewInternalServerError() *AppError {
	return NewAppError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(http.StatusNotFound, message)
}

func NewConflictError(message string) *AppError {
	return NewAppError(http.StatusConflict, message)
}

func NewValidationFailedError(message string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(http.StatusForbidden, message)
}
