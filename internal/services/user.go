package services

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/entities"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/errs"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
	"context"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type UserService interface {
	GetUserState(ctx context.Context, game, integrator string, lordParams interface{}) (*entities.SessionState, error)
	GetUserStateBySessionToken(ctx context.Context, sessionToken string) (*entities.SessionState, error)
	Wager(ctx context.Context, sessionToken, freeSpinID, currency string, wager int64) (*entities.SessionState, error)
}

type Service struct {
	core       *CoreService
	historySrv *HistoryService
	lordClient overlord.Client
}

func NewUserService(core *CoreService, historySrv *HistoryService, lordClient overlord.Client) UserService {
	return &Service{
		core:       core,
		historySrv: historySrv,
		lordClient: lordClient,
	}
}

func (s *Service) GetUserState(ctx context.Context, game, integrator string, lordParams interface{}) (*entities.SessionState, error) {
	state := &entities.SessionState{}

	stateOut, err := s.lordClient.InitUserState(ctx, game, integrator, lordParams)
	if err != nil {
		zap.S().Error("failed to init betlord user state", err)

		return nil, errs.TranslateOverlordErr(err)
	}

	user, err := entities.UserStateFromOverlordState(stateOut)
	if err != nil {
		return nil, err
	}

	state.FillByUserState(user)

	history, err := s.historySrv.Get(ctx, user.UserID, user.Game)
	if err == nil {
		if lo.Contains(state.WagerLevels, history.Bet) {
			state.DefaultWager = history.Bet
		}

		state.Reels = history.ReelsState
		state.SpinsIndexes.BaseSpinIndex = history.SpinIndexes.BaseSpinIndex
		state.SpinsIndexes.BonusSpinIndex = history.SpinIndexes.BonusSpinIndex
	}

	return state, nil
}

func (s *Service) GetUserStateBySessionToken(ctx context.Context, sessionToken string) (*entities.SessionState, error) {
	state := &entities.SessionState{}

	stateOut, err := s.lordClient.GetStateBySessionToken(ctx, sessionToken)
	if err != nil {
		zap.S().Error("failed to get user state by session token", err)

		return nil, errs.TranslateOverlordErr(err)
	}

	user, err := entities.UserStateFromOverlordState(stateOut)
	if err != nil {
		return nil, err
	}

	state.FillByUserState(user)

	return state, nil
}

func (s *Service) Wager(ctx context.Context, sessionToken, freeSpinID, currency string, wager int64) (*entities.SessionState, error) {
	state, err := s.GetUserStateBySessionToken(ctx, sessionToken)
	if err != nil {
		return nil, err
	}

	if freeSpinID != "" {
		fs, err := s.findUserFreeSpin(ctx, sessionToken, freeSpinID)
		if err != nil {
			return nil, err
		}

		wager = int64(fs.Value)
	}

	if freeSpinID == "" && state.Balance < wager {
		return nil, errs.ErrNotEnoughMoney
	}

	if freeSpinID == "" && !lo.Contains(state.WagerLevels, wager) {
		return nil, errs.NewInternalValidationErrorFromString(errs.OneOfListError("wager", state.WagerLevels))
	}

	user := s.core.BaseUser()

	user.State.UserID = state.UserID
	user.State.ExternalUserID = state.ExternalUserID
	user.State.Game = state.Game
	user.State.GameID = state.GameID
	user.State.SessionToken = state.SessionToken
	user.State.FreespinID = freeSpinID
	user.State.RoundID = uuid.New()
	user.State.Integrator = state.Integrator
	user.State.Operator = state.Operator
	user.State.Provider = state.Provider
	user.State.Currency = state.Currency
	user.State.StartBalance = state.Balance
	user.State.WagerAmount = wager

	err = s.core.ExecSpin(user)
	if err != nil {
		return nil, err
	}

	res, err := s.lordClient.
		AtomicBet(ctx, user.State.SessionToken, user.State.FreespinID, user.State.RoundID.String(), wager, user.State.TotalWins)
	if err != nil {
		zap.S().Error("failed to bet", err)

		return nil, errs.TranslateOverlordErr(err)
	}

	user.State.Balance = res.Balance
	user.State.TransactionID = res.TransactionId

	err = s.historySrv.Create(ctx, user)
	if err != nil {
		zap.S().Error("failed to save spin", err)

		return nil, err
	}

	return &user.State, nil
}

func (s *Service) findUserFreeSpin(ctx context.Context, session, fsID string) (*entities.FreeSpin, error) {
	pureFS, err := s.lordClient.GetAvailableFreeSpins(ctx, session)
	if err != nil {
		return nil, errs.TranslateOverlordErr(err)
	}

	fs := entities.FreeSpinsFromLord(pureFS.FreeBets)
	item, ok := lo.Find(fs, func(item *entities.FreeSpin) bool {
		return item.ID == fsID
	})

	if !ok {
		return nil, errs.ErrWrongFreeSpinID
	}

	return item, nil
}
