package custom_errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorCode string

const (
	ErrStorageConnection ErrorCode = "err_storage_connection"
)

type CustomError struct {
	Code ErrorCode
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%s", e.Code)
}

// Optionnel : constructeur
func New(code ErrorCode) error {
	return errors.WithStack(CustomError{Code: code})
}
