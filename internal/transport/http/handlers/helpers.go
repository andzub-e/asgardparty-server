package handlers

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/entities"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/errs"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

var errorMap = map[error]func(ctx *gin.Context, data interface{}, meta map[string]interface{}){
	errs.ErrNotEnoughMoney:           http.PaymentRequired,
	errs.ErrBadDataGiven:             http.BadRequest,
	errs.ErrWrongSessionToken:        http.Unauthorized,
	errs.ErrSessionTokenExpired:      http.SessionExpired,
	errs.ErrWrongFreeSpinID:          http.Conflict,
	errs.ErrHistoryNotFound:          http.Conflict,
	errs.ErrUserHasDifferentCurrency: http.Conflict,

	errs.ErrUserIsBlocked:             http.Forbidden,
	errs.ErrIntegratorCriticalFailure: http.ServiceUnavailableError,
}

func handleServiceError(ctx *gin.Context, err error) {
	internalValidationError, ok := err.(errs.InternalValidationError)
	if ok {
		http.ValidationFailed(ctx, internalValidationError)

		return
	}

	fn, ok := errorMap[err]
	if !ok {
		http.ServerError(ctx, err, nil)

		return
	}

	fn(ctx, err, nil)
}

func patchContextData(ctx *gin.Context, parsedRequest interface{}) error {
	requestBody, err := json.Marshal(parsedRequest)

	if err != nil {
		return err
	}

	ctx.Request = ctx.Request.WithContext(
		context.WithValue(ctx.Request.Context(), constants.CtxAdditionalDataKey,
			entities.NewPlayerMetaData(ctx.ClientIP(),
				ctx.GetHeader("User-Agent"),
				ctx.Request.Header.Get("Origin"),
				requestBody)),
	)

	return nil
}
