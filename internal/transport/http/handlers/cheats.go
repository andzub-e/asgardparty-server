package handlers

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/services"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/requests"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/responses"
	"github.com/gin-gonic/gin"
)

var cheatsSuccessResponse = responses.NoContentResponse{Success: true}

type cheatsHandler struct {
	cheatsService     *services.CheatsService
	isCheatsAvailable bool
}

func NewCheatsHandler(cheatsService *services.CheatsService, isCheatsAvailable bool) http.Handler {
	return &cheatsHandler{
		cheatsService:     cheatsService,
		isCheatsAvailable: isCheatsAvailable,
	}
}

func (h *cheatsHandler) Register(router *gin.RouterGroup) {
	if h.isCheatsAvailable {
		router.POST("/cheat/custom_figures", h.customFigures)
	}
}

// @Summary Cheat custom figures
// @Tags cheat
// @Consume application/json
// @Description Set next spin to be with custom figures from request
// @Param JSON body requests.CheatCustomFiguresRequest true "body for exec cheat"
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.NoContentResponse
// @Failure 400 {object} http.Response
// @Failure 500 {object} http.Response
// @Router /cheat/custom_figures [post]
func (h *cheatsHandler) customFigures(ctx *gin.Context) {
	var req requests.CheatCustomFiguresRequest

	if err := ctx.ShouldBind(&req); err != nil {
		http.BadRequest(ctx, err, nil)

		return
	}

	err := h.cheatsService.CustomFigures(req.SessionToken, req.Figures)

	if err != nil {
		http.BadRequest(ctx, err, nil)

		return
	}

	http.OK(ctx, cheatsSuccessResponse, nil)
}
