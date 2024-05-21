package handlers

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/responses"
	"github.com/gin-gonic/gin"
)

var tag = "no tag"

type metaHandler struct{}

func NewMetaHandler() http.Handler {
	return &metaHandler{}
}

func (h *metaHandler) Register(route *gin.RouterGroup) {
	route.GET("health", h.health)
	route.GET("info", h.info)
}

// @Summary Check health.
// @Tags meta
// @Consume application/json
// @Description Check service health.
// @Accept  json
// @Produce  json
// @Success 200  {object} responses.HealthResponse
// @Router /health [get]
func (h *metaHandler) health(ctx *gin.Context) {
	http.OK(ctx, responses.HealthResponse{Success: "ok"}, nil)
}

// @Summary Check tag.
// @Tags meta
// @Consume application/json
// @Description Check service tag.
// @Accept  json
// @Produce  json
// @Success 200  {object} responses.InfoResponse
// @Router /info [get]
func (h *metaHandler) info(ctx *gin.Context) {
	http.OK(ctx, responses.InfoResponse{Tag: tag}, nil)
}
