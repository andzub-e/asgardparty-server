package handlers

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/services"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/requests"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/responses"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var noContentSuccessResponse = responses.NoContentResponse{Success: true}

type coreHandler struct {
	userService       services.UserService
	freeSpinsSrv      *services.FreeSpinService
	historyService    *services.HistoryService
	maxProcessingTime time.Duration
}

func NewCoreHandler(userService services.UserService, historyService *services.HistoryService,
	freeSpinsSrv *services.FreeSpinService, maxProcessingTime time.Duration) http.Handler {
	return &coreHandler{
		userService:       userService,
		historyService:    historyService,
		freeSpinsSrv:      freeSpinsSrv,
		maxProcessingTime: maxProcessingTime,
	}
}

func (h *coreHandler) Register(router *gin.RouterGroup) {
	core := router.Group("core")
	freespins := core.Group("free_spins")

	core.POST("/state", h.getUserState)
	core.POST("/wager", h.placeWager)

	core.GET("/spins_history", h.getSpinsHistory)
	core.POST("/spin_indexes/update", h.updateSpinIndexes)

	freespins.GET("get", h.getFreeSpins)
	freespins.GET("decline", h.declineFreeSpins)

	freespins.GET("get_with_integrator_bet", h.getFreeSpinsWithIntegratorBet)
	freespins.GET("cancel_with_integrator_bet", h.cancelFreeSpinsWithIntegratorBet)
}

// @Summary State
// @Tags core
// @Consume application/json
// @Param session_id query string true "session id from game client"
// @Param integrator query string true "integrator name"
// @Param game query string true "game id"
// @Param params query string true "special integrator parameters"
// @Description Retrieves an initial state of the game from the bet-overlord service.
// @Description For every new session you need to send /state request. /br In response you'll get a session token for a bets placing.
// @Description Mock data: integrator - MOCK, game - TEST, params example - { "user": "58e361be-2edc-4b4e-bf24-5a348a5eff3c", "token": "04e717e7-e5af-42eb-8c67-08aa647b5c7b", "game": "test", "integrator": "MOCK" }
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.StateResponse
// @Failure 400 {object} http.Response
// @Failure 500 {object} http.Response
// @Router /core/state [post]
func (h *coreHandler) getUserState(ctx *gin.Context) {
	var req requests.StateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	userState, err := h.userService.GetUserState(ctxWithTimeout, req.Game, req.Integrator, req.Params)
	if err != nil {
		zap.S().Error("failed to get user state from db", err)
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, responses.StateResponse{SessionState: *userState}, nil)
}

// @Summary Place a bet
// @Tags core
// @Consume application/json
// @Param JSON body requests.PlaceWagerRequest true "body for a making bet(spin)"
// @Description Make a bet (spin).
// @Accept  json
// @Produce  json
// @Success 200 {object} entities.SessionState
// @Failure 400 {object} http.Response
// @Failure 500 {object} http.Response
// @Router /core/wager [post]
func (h *coreHandler) placeWager(ctx *gin.Context) {
	var req requests.PlaceWagerRequest
	if err := ctx.ShouldBind(&req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	sessionState, err := h.userService.Wager(ctxWithTimeout, req.SessionToken, req.FreeSpinID, req.Currency, req.Wager)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, sessionState, nil)
}

// @Summary Get spins history
// @Tags core
// @Consume application/json
// @Description Retrieves user's spins history
// @Param session_token query string true "session token"
// @Param page query int true "page"
// @Param count query int true "count"
// @Produce  json
// @Success 200 {object} entities.HistoryPagination
// @Failure 400 {object} http.Response
// @Failure 500 {object} http.Response
// @Router /core/spins_history [get]
func (h *coreHandler) getSpinsHistory(ctx *gin.Context) {
	req := &requests.HistoryRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	state, err := h.userService.GetUserStateBySessionToken(ctxWithTimeout, req.SessionToken)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	historyPagination, err := h.historyService.HistoryPagination(ctxWithTimeout, state.UserID, state.Game, *req.Count, *req.Page)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, historyPagination, nil)
}

// @Tags core
// @Description get available free spins
// @Accept json
// @Produce json
// @Param {object} session_token query requests.GetFreeSpinsRequest true "get free spins"
// @Success 200 {object} responses.GetFreeSpinsResponse
// @Failure 400 {object} http.Response
// @Failure 404 {object} http.Response
// @Failure 500 {object} http.Response
// @Router /core/free_spins/get [get]
func (h *coreHandler) getFreeSpins(ctx *gin.Context) {
	var req requests.GetFreeSpinsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	freeSpins, err := h.freeSpinsSrv.GetFreeSpins(ctxWithTimeout, req.SessionToken)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, responses.GetFreeSpinsResponse{FreeSpins: freeSpins}, nil)
}

func (h *coreHandler) declineFreeSpins(ctx *gin.Context) {
	req := &requests.GetFreeSpinsRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	if err := h.freeSpinsSrv.CancelFreeSpins(ctxWithTimeout, req.SessionToken); err != nil {
		zap.S().Error("failed to decline free spins", err)
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, noContentSuccessResponse, nil)
}

// @Tags core
// @Description important restoring endpoint which gives an opportunity to track shown spins
// @Accept json
// @Produce json
// @Param {object} spin_indexes body requests.UpdateSpinsIndexesRequest true "update spin indexes"
// @Success 200 {object} responses.NoContentResponse
// @Failure 400 {object} http.Response
// @Failure 404 {object} http.Response
// @Failure 500 {object} http.Response
// @Router /core/spin_indexes/update [post]
func (h *coreHandler) updateSpinIndexes(ctx *gin.Context) {
	req := requests.UpdateSpinsIndexesRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	state, err := h.userService.GetUserStateBySessionToken(ctxWithTimeout, req.SessionToken)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	err = h.historyService.UpdateLastSpinsIndexes(ctxWithTimeout, state.UserID, state.Game, *req.BaseSpinIndex, *req.BonusSpinIndex)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, noContentSuccessResponse, nil)
}

func (h *coreHandler) getFreeSpinsWithIntegratorBet(ctx *gin.Context) {
	var req requests.GetFreeSpinsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	freeSpins, err := h.freeSpinsSrv.GetFreeSpinsWithIntegratorBet(ctxWithTimeout, req.SessionToken)
	if err != nil {
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, responses.GetFreeSpinsWithIntegratorBetResponse{FreeSpins: freeSpins}, nil)
}

func (h *coreHandler) cancelFreeSpinsWithIntegratorBet(ctx *gin.Context) {
	var req requests.FreeSpinsWithIntegratorBetRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.ValidationFailed(ctx, err)

		return
	}

	if err := patchContextData(ctx, req); err != nil {
		zap.S().Error("http: failed to patch context", err)
		http.BadRequest(ctx, err, nil)

		return
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx.Request.Context(), h.maxProcessingTime)
	defer cancel()

	if err := h.freeSpinsSrv.CancelFreeSpinsWithIntegratorBet(ctxWithTimeout, req.SessionToken, req.IntegratorBetId); err != nil {
		zap.S().Error("cancel fs: ", err)
		handleServiceError(ctx, err)

		return
	}

	http.OK(ctx, noContentSuccessResponse, nil)
}
