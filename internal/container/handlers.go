package container

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/config"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/services"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/transport/http/handlers"
	"github.com/sarulabs/di"
)

func BuildHandlers() []di.Def {
	return []di.Def{
		{
			Name: constants.CheatsHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cheatsService := ctn.Get(constants.CheatsServiceName).(*services.CheatsService)
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				return handlers.NewCheatsHandler(cheatsService, cfg.EngineConfig.IsCheatsAvailable), nil
			},
		},
		{
			Name: constants.CoreHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get(constants.ConfigName).(*config.Config)

				userService := ctn.Get(constants.UserServiceName).(services.UserService)
				historyService := ctn.Get(constants.HistoryServiceName).(*services.HistoryService)
				freeSpinService := ctn.Get(constants.FreeSpinsServiceName).(*services.FreeSpinService)

				return handlers.NewCoreHandler(userService, historyService, freeSpinService, cfg.ServerConfig.MaxProcessingTime), nil
			},
		},
		{
			Name: constants.MetaHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				return handlers.NewMetaHandler(), nil
			},
		},
		{
			Name: constants.MetricsHandlerName,
			Build: func(ctn di.Container) (interface{}, error) {
				return handlers.NewMetricsHandler()
			},
		},
	}
}
