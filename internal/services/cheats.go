package services

import (
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
	"go.uber.org/zap"
)

type CheatsService struct {
	core *CoreService
}

func NewCheatsService(core *CoreService) *CheatsService {
	return &CheatsService{
		core: core,
	}
}

func (s *CheatsService) CustomFigures(sessionToken string, figures []string) error {
	data, err := s.core.CheatCustomFigures(figures)
	if err != nil {
		zap.S().Error("stops validation failed", err)

		return err
	}

	err = utils.Cache.Set(sessionToken+"-base", data)
	if err != nil {
		zap.S().Error("failed to add data to cache", err)

		return err
	}

	return nil
}
