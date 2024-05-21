package container

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/config"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/validator"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/history"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/rng"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/tracer"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"sync"
)

var container di.Container
var once sync.Once

func Build(ctx context.Context, wg *sync.WaitGroup) di.Container {
	once.Do(func() {
		builder, _ := di.NewBuilder()
		defs := []di.Def{
			{
				Name: constants.LoggerName,
				Build: func(ctn di.Container) (interface{}, error) {
					logger, err := zap.NewProduction()
					if constants.Debug {
						logger, err = zap.NewDevelopment()
					}

					if err != nil {
						return nil, fmt.Errorf("can't initialize zap logger: %v", err)
					}

					zap.ReplaceGlobals(logger)

					return logger, nil
				},
			},
			{
				Name: constants.ConfigName,
				Build: func(ctn di.Container) (interface{}, error) {
					return config.New()
				},
			},
			{
				Name: constants.TracerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return tracer.NewTracer(cfg.TracerConfig)
				},
			},
			{
				Name: constants.OverlordName,
				Build: func(ctn di.Container) (interface{}, error) {
					conf := ctn.Get(constants.ConfigName).(*config.Config)

					return overlord.NewClient(conf.OverlordConfig)
				},
			},
			{
				Name: constants.HistoryName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)

					return history.NewClient(cfg.HistoryConfig)
				},
			},
			{
				Name: constants.RNGName,
				Build: func(ctn di.Container) (interface{}, error) {
					conf := ctn.Get(constants.ConfigName).(*config.Config)

					return rng.New(conf.RNGConfig)
				},
			},
			{
				Name: constants.ValidatorName,
				Build: func(ctn di.Container) (interface{}, error) {
					conf := ctn.Get(constants.ConfigName).(*config.Config)

					return validator.New(conf.ConstantsConfig)
				},
			},
			{
				Name: constants.ServerName,
				Build: func(ctn di.Container) (interface{}, error) {
					cfg := ctn.Get(constants.ConfigName).(*config.Config)
					handlers := []http.Handler{
						ctn.Get(constants.CoreHandlerName).(http.Handler),
						ctn.Get(constants.CheatsHandlerName).(http.Handler),
						ctn.Get(constants.MetaHandlerName).(http.Handler),
						ctn.Get(constants.MetricsHandlerName).(http.Handler),
					}

					var middlewares = []func(ctx *gin.Context){
						ctn.Get(constants.HTTPCorsMiddlewareName).(func(ctx *gin.Context)),
						ctn.Get(constants.HTTPLogMiddlewareName).(func(ctx *gin.Context)),
						ctn.Get(constants.HTTPTraceMiddlewareName).(func(ctx *gin.Context)),
					}

					return http.New(ctx, wg, cfg.ServerConfig, handlers, middlewares), nil
				},
			},
		}

		defs = append(defs, BuildServices()...)
		defs = append(defs, BuildHandlers()...)
		defs = append(defs, BuildMiddlewares()...)

		if err := builder.Add(defs...); err != nil {
			panic(err)
		}

		container = builder.Build()
	})

	return container
}
