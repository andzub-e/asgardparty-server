package errs

import (
	"errors"
	"fmt"
)

type InternalValidationError struct {
	Err error
}

func (err InternalValidationError) Error() string {
	return err.Err.Error()
}

func NewInternalValidationError(err error) InternalValidationError {
	return InternalValidationError{Err: err}
}

func NewInternalValidationErrorFromString(str string) InternalValidationError {
	return NewInternalValidationError(errors.New(str))
}

var (
	ErrNotEnoughMoney           = errors.New("not enough money")
	ErrWrongSessionToken        = errors.New("wrong session token")
	ErrSessionTokenExpired      = errors.New("session token expired")
	ErrUserHasDifferentCurrency = errors.New("user_has_different_currency")
	ErrWrongFreeSpinID          = errors.New("wrong free spin id")
	ErrBadDataGiven             = errors.New("bad data given")
	ErrHistoryNotFound          = errors.New("history not found")

	ErrUserIsBlocked             = errors.New("user is blocked")
	ErrIntegratorCriticalFailure = errors.New("integrator critical failure")
)

func OneOfListError[T comparable](field string, list []T) string {
	return fmt.Sprintf("%s must be one of %v", field, list)
}
