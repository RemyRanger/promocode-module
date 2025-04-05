package custom_errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type ErrorCode string

const (
	ErrPromocodeExist ErrorCode = "err_promocode_exist"
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
