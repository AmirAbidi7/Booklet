package errors

import "fmt"

type ErrorKind string

const (
	// Auth errors
	ErrorKindUserNotFound    ErrorKind = "USER_NOT_FOUND"
	ErrorKindUserExists      ErrorKind = "USER_EXISTS"
	ErrorKindInvalidPassword ErrorKind = "INVALID_PASSWORD"
	ErrorKindUnauthorized    ErrorKind = "UNAUTHORIZED"
	// PDF errors
	ErrorKindPdfNotFound  ErrorKind = "PDF_NOT_FOUND"
	ErrorKindInvalidPdf   ErrorKind = "INVALID_PDF"
	ErrorKindFileTooLarge ErrorKind = "FILE_TOO_LARGE"
	// Database errors
	ErrorKindDatabase ErrorKind = "DATABASE_ERROR"
	// Storage errors
	ErrorKindStorageError ErrorKind = "STORAGE_ERROR"
	// Validation errors
	ErrorKindValidation ErrorKind = "VALIDATION_ERROR"
	// Internal server errors
	ErrorKindInternal ErrorKind = "INTERNAL_ERROR"
)

type AppError struct {
	Kind       ErrorKind
	Message    string
	Err        error
	StatusCode int
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Kind, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func UserExists(email string) *AppError {
	return &AppError{
		Kind:       ErrorKindUserExists,
		Message:    fmt.Sprintf("user with email %s already exists", email),
		StatusCode: 409,
	}
}

func UserNotFound() *AppError {
	return &AppError{
		Kind:       ErrorKindUserNotFound,
		Message:    "user not found",
		StatusCode: 404,
	}
}

func InvalidPassword() *AppError {
	return &AppError{
		Kind:       ErrorKindInvalidPassword,
		Message:    "invalid email or password",
		StatusCode: 401,
	}
}

func PdfNotFound(id uint) *AppError {
	return &AppError{
		Kind:       ErrorKindPdfNotFound,
		Message:    fmt.Sprintf("pdf with id %d not found", id),
		StatusCode: 404,
	}
}

func FileTooLarge(maxMB int) *AppError {
	return &AppError{
		Kind:       ErrorKindFileTooLarge,
		Message:    fmt.Sprintf("file size cannot exceed %dMB", maxMB),
		StatusCode: 400,
	}
}

func InvalidContentType(contentType string) *AppError {
	return &AppError{
		Kind:       ErrorKindInvalidPdf,
		Message:    fmt.Sprintf("invalid content type: %s. Expected application/pdf", contentType),
		StatusCode: 400,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Kind:       ErrorKindUnauthorized,
		Message:    message,
		StatusCode: 401,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Kind:       ErrorKindUnauthorized,
		Message:    message,
		StatusCode: 403,
	}
}

func ValidationError(message string) *AppError {
	return &AppError{
		Kind:       ErrorKindValidation,
		Message:    message,
		StatusCode: 422,
	}
}

func Database(err error, message string) *AppError {
	return &AppError{
		Kind:       ErrorKindDatabase,
		Message:    message,
		Err:        err,
		StatusCode: 500,
	}
}

func StorageError(err error, message string) *AppError {
	return &AppError{
		Kind:       ErrorKindStorageError,
		Message:    message,
		Err:        err,
		StatusCode: 500,
	}
}

func Internal(err error, message string) *AppError {
	return &AppError{
		Kind:       ErrorKindInternal,
		Message:    message,
		Err:        err,
		StatusCode: 500,
	}
}

// Helper to check error type
func Is(err error, kind ErrorKind) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Kind == kind
	}
	return false
}
