package http

import (
	"context"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"sync"
	"time"

	"bitbucket.org/electronicjaw/asgardparty-server/docs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	wg     *sync.WaitGroup
	ctx    context.Context
	server *http.Server
	router *gin.Engine
}

func New(ctx context.Context, wg *sync.WaitGroup, cfg *Config, handlers []Handler, middlewares []func(ctx *gin.Context)) *Server {
	docs.SwaggerInfo.Title = "API"
	docs.SwaggerInfo.Description = "This is an Asgard Party slot server."
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	s := &Server{
		wg:  wg,
		ctx: ctx,
		server: &http.Server{
			Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler:           nil,
			ReadHeaderTimeout: 10 * time.Second,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       30 * time.Second,
		},
		router: gin.New(),
	}

	s.registerMiddlewares(middlewares)

	s.router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Asgard Party API")
	})
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := s.router.Group("")

	s.registerHandlers(api, handlers...)

	return s
}

func (s *Server) registerMiddlewares(middlewares []func(ctx *gin.Context)) {
	for _, mw := range middlewares {
		s.router.Use(mw)
	}
}

func (s *Server) registerHandlers(api *gin.RouterGroup, handlers ...Handler) {
	for _, h := range handlers {
		h.Register(api)
	}

	s.server.Handler = s.router
}

func (s *Server) Run() {
	s.wg.Add(1)
	zap.S().Infof("server listening: %s", s.server.Addr)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.S().Error(err.Error())
	}
}

func (s *Server) Shutdown() error {
	zap.S().Info("Shutdown server...")
	zap.S().Info("Stopping http server...")

	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)

	defer func() {
		cancel()
		s.wg.Done()
	}()

	if err := s.server.Shutdown(ctx); err != nil {
		zap.S().Fatal("Server forced to shutdown:", zap.Error(err))

		return err
	}

	zap.S().Info("Server successfully stopped.")

	return nil
}
