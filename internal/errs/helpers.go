package errs

import (
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/history"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
)

var translateMap = map[error]error{
	overlord.ErrMarshaling:                ErrBadDataGiven,
	overlord.ErrBalanceTooLow:             ErrNotEnoughMoney,
	overlord.ErrWrongSessionToken:         ErrWrongSessionToken,
	overlord.ErrWrongFreeSpinID:           ErrWrongFreeSpinID,
	overlord.ErrSessionTokenExpired:       ErrSessionTokenExpired,
	overlord.ErrUserHasDifferentCurrency:  ErrUserHasDifferentCurrency,
	overlord.ErrUserIsBlocked:             ErrUserIsBlocked,
	overlord.ErrIntegratorCriticalFailure: ErrIntegratorCriticalFailure,
}

var translateHistoryMap = map[error]error{
	history.ErrSpinNotFound: ErrHistoryNotFound,
}

func TranslateOverlordErr(err error) error {
	validationErr, ok := err.(overlord.ValidationError)
	if ok {
		return InternalValidationError{Err: validationErr}
	}

	res, ok := translateMap[err]
	if !ok {
		return err
	}

	return res
}

func TranslateHistoryErr(err error) error {
	validationErr, ok := err.(overlord.ValidationError)
	if ok {
		return InternalValidationError{Err: validationErr}
	}

	res, ok := translateHistoryMap[err]
	if !ok {
		return err
	}

	return res
}
