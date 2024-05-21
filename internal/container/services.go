package container

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/services"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/history"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
	"github.com/sarulabs/di"
)

func BuildServices() []di.Def {
	return []di.Def{
		{
			Name: constants.CoreServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				return services.NewCoreService(), nil
			},
		},
		{
			Name: constants.UserServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				coreService := ctn.Get(constants.CoreServiceName).(*services.CoreService)
				historySrv := ctn.Get(constants.HistoryServiceName).(*services.HistoryService)
				overlordClient := ctn.Get(constants.OverlordName).(overlord.Client)

				return services.NewUserService(coreService, historySrv, overlordClient), nil
			},
		},
		{
			Name: constants.CheatsServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				coreService := ctn.Get(constants.CoreServiceName).(*services.CoreService)

				return services.NewCheatsService(coreService), nil
			},
		},
		{
			Name: constants.HistoryServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				historyClient := ctn.Get(constants.HistoryName).(history.Client)

				return services.NewHistoryService(historyClient), nil
			},
		},
		{
			Name: constants.FreeSpinsServiceName,
			Build: func(ctn di.Container) (interface{}, error) {
				overlordClient := ctn.Get(constants.OverlordName).(overlord.Client)

				return services.NewFreeSpinService(overlordClient), nil
			},
		},
	}
}
