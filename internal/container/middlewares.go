package container

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/middlewares"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/tracer"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
	"time"
)

func BuildMiddlewares() []di.Def {
	return []di.Def{
		{
			Name: constants.HTTPCorsMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				return middlewares.CORSMiddleware(), nil
			},
		},
		{
			Name: constants.HTTPLogMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := ctn.Get(constants.LoggerName).(*zap.Logger)

				return middlewares.Log(logger.Sugar(), time.RFC3339, true), nil
			},
		},
		{
			Name: constants.HTTPTraceMiddlewareName,
			Build: func(ctn di.Container) (interface{}, error) {
				tr := ctn.Get(constants.TracerName).(*tracer.JaegerTracer)

				return middlewares.TraceMiddleware(tr), nil
			},
		},
	}
}
