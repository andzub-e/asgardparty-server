package services

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/constants"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/entities"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/errs"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/history"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HistoryService struct {
	historyClient history.Client
}

func NewHistoryService(historyClient history.Client) *HistoryService {
	return &HistoryService{
		historyClient: historyClient,
	}
}

func (s *HistoryService) HistoryPagination(ctx context.Context, userID uuid.UUID, game string, count, page int) (entities.HistoryPagination, error) {
	pagination, err := s.historyClient.Pagination(ctx, userID, game, count, page)
	if err != nil {
		return entities.HistoryPagination{}, err
	}

	p := entities.HistoryPagination{
		Spins:       []*entities.History{},
		CurrentPage: int(pagination.Page),
		Count:       int(pagination.Limit),
		Total:       int(pagination.Total)}

	for _, item := range pagination.Items {
		r, err := entities.SpinOutToHistory(item)

		if err != nil {
			return entities.HistoryPagination{}, err
		}

		p.Spins = append(p.Spins, r)
	}

	return p, nil
}

func (s *HistoryService) Get(ctx context.Context, userID uuid.UUID, game string) (*entities.History, error) {
	zap.S().Info(s.historyClient)

	spinOut, err := s.historyClient.LastRecord(ctx, userID, game)
	if err != nil {
		return nil, errs.TranslateHistoryErr(err)
	}

	return entities.SpinOutToHistory(spinOut)
}

func (s *HistoryService) UpdateLastSpinsIndexes(ctx context.Context, userID uuid.UUID, gameName string, baseSpinIndex, bonusSpinIndex int) error {
	hr, err := s.Get(ctx, userID, gameName)
	if err != nil {
		return err
	}

	hr.SpinIndexes.BaseSpinIndex = baseSpinIndex
	hr.SpinIndexes.BonusSpinIndex = bonusSpinIndex
	hr.IsShown = true

	metaData := s.extractAdditionalData(ctx)

	in, err := hr.ToHistoryIn(metaData)
	if err != nil {
		return err
	}

	return s.historyClient.Update(ctx, in)
}

func (s *HistoryService) Create(ctx context.Context, user *entities.User) error {
	record, err := user.ExtractHistory()
	if err != nil {
		return nil
	}

	metaData := s.extractAdditionalData(ctx)

	in, err := record.ToHistoryIn(metaData)
	if err != nil {
		return err
	}

	return s.historyClient.Create(ctx, in)
}

func (s *HistoryService) extractAdditionalData(ctx context.Context) *entities.PlayerMetaData {
	val, ok := ctx.Value(constants.CtxAdditionalDataKey).(*entities.PlayerMetaData)
	if !ok {
		zap.S().Warn("there is no additional data")
	}

	return val
}
