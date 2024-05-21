package handlers

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metricsHandler struct{}

func NewMetricsHandler() (http.Handler, error) {
	return &metricsHandler{}, nil
}

func (m metricsHandler) Register(router *gin.RouterGroup) {
	metrics := router.Group("metrics")

	fn := promhttp.Handler()

	metrics.GET("/", func(ginCtx *gin.Context) {
		fn.ServeHTTP(ginCtx.Writer, ginCtx.Request)
	})
}
