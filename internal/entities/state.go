package entities

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/errs"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/game"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
	"github.com/google/uuid"
	"time"
)

type SessionState struct {
	Username     string       `json:"username,omitempty" swaggertype:"string"`
	Balance      int64        `json:"balance" swaggertype:"integer"`
	Currency     string       `json:"currency,omitempty" swaggertype:"integer"`
	WagerID      string       `json:"wallet_play_id,omitempty" swaggertype:"string"`
	WagerAmount  int64        `json:"last_wager,omitempty" swaggertype:"integer"`
	FreespinID   string       `json:"freespinid,omitempty" swaggertype:"integer"`
	Reels        game.State   `json:"reels"`
	TotalWins    int64        `json:"total_wins" swaggertype:"integer"`
	SessionToken string       `json:"session_token" swaggertype:"string"`
	WagerLevels  []int64      `json:"wager_levels,omitempty"`
	DefaultWager int64        `json:"default_wager,omitempty"`
	Game         string       `json:"THE_EJAW_SLOT,omitempty"`
	GameID       uuid.UUID    `json:"game_id,omitempty"`
	SpinsIndexes SpinsIndexes `json:"spins_indexes"`

	StartBalance   int64     `json:"-"`
	UserID         uuid.UUID `json:"-"`
	ExternalUserID string    `json:"-"`
	Integrator     string    `json:"-"`
	Operator       string    `json:"-"`
	Provider       string    `json:"-"`
	TransactionID  string    `json:"-"`
	RoundID        uuid.UUID `json:"-"`

	IsDemo bool `json:"is_demo"`

	Error string `json:"error,omitempty"`
}

type User struct {
	State SessionState
}

type UserState struct {
	UserID         uuid.UUID `json:"user_id" example:"some_id1"`
	ExternalUserID string    `json:"external_user_id" example:"some_id1"`
	Integrator     string    `json:"integrator,omitempty" example:"mock"`
	Operator       string    `json:"Operator,omitempty" example:"mock"`
	Provider       string    `json:"provider,omitempty" example:"mock"`
	Game           string    `json:"THE_EJAW_SLOT" example:"test"`
	GameID         uuid.UUID `json:"game_id" example:"test"`
	Username       string    `json:"username" example:"mock_state"`
	SessionToken   string    `json:"session_token" example:"9eedb051-32f6-4e24-9caa-93fa63685a95"`
	Balance        int64     `json:"balance" example:"9997170"`
	Currency       string    `json:"currency" example:"XTS"`

	DefaultWager       int64   `json:"default_wager" example:"2000"`
	CurrencyMultiplier int64   `json:"currency_multiplier"`
	WagerLevels        []int64 `json:"wager_levels" example:"100,200,500,700,1000,1500,2000,2500,5000,10000,25000,50000,100000,150000,250000,500000"`

	IsDemo bool `json:"is_demo"`
}

// SpinsIndexes BaseSpinIndex can be 0 or 1 (if there are no win combination - it will be always zero)
// if spin has a win combination 0 means that it's not shown, accordingly 1 means spin have been shown
// SpinsIndexes BonusSpinIndex - can be form 0 to infinity (max bonus spin quantity),
// and that field shows how much of them have been shown.
type SpinsIndexes struct {
	BaseSpinIndex  int `json:"base_spin_index"`
	BonusSpinIndex int `json:"bonus_spin_index"`
}

func (state *SessionState) FillByUserState(user *UserState) {
	state.Username = user.Username
	state.Balance = user.Balance
	state.Currency = user.Currency
	state.SessionToken = user.SessionToken
	state.WagerLevels = user.WagerLevels
	state.DefaultWager = user.DefaultWager
	state.Game = user.Game
	state.GameID = user.GameID
	state.UserID = user.UserID
	state.ExternalUserID = user.ExternalUserID
	state.Integrator = user.Integrator
	state.Operator = user.Operator
	state.Provider = user.Provider
}

func UserStateFromOverlordState(state *overlord.InitUserStateOut) (*UserState, error) {
	userID, err := uuid.Parse(state.UserId)
	if err != nil {
		return nil, err
	}

	gameID, err := uuid.Parse(state.GameId)
	if err != nil {
		return nil, err
	}

	u := &UserState{}
	u.UserID = userID
	u.ExternalUserID = state.ExternalUserId
	u.Integrator = state.Integrator
	u.Operator = state.Operator
	u.Provider = state.Provider
	u.Game = state.Game
	u.GameID = gameID
	u.Username = state.Username
	u.SessionToken = state.SessionToken
	u.Balance = state.Balance
	u.Currency = state.Currency

	u.DefaultWager = state.DefaultWager
	u.CurrencyMultiplier = state.CurrencyMultiplier
	u.WagerLevels = state.WagerLevels

	u.IsDemo = state.IsDemo

	return u, nil
}

func (u *User) ExtractHistory() (*History, error) {
	transactionID, err := uuid.Parse(u.State.TransactionID)
	if err != nil {
		// rollback
		return nil, errs.ErrBadDataGiven
	}

	session, err := uuid.Parse(u.State.SessionToken)
	if err != nil {
		return nil, err
	}

	return &History{
		UserID:         u.State.UserID,
		ExternalUserID: u.State.ExternalUserID,
		Game:           u.State.Game,
		GameID:         u.State.GameID,
		TransactionID:  transactionID,
		SessionToken:   session,
		Integrator:     u.State.Integrator,
		Operator:       u.State.Operator,
		Provider:       u.State.Provider,

		Currency:     u.State.Currency,
		StartBalance: u.State.StartBalance,
		Bet:          u.State.WagerAmount,
		EndBalance:   u.State.Balance,
		ID:           u.State.RoundID,
		ReelsState:   u.State.Reels,
		StartTime:    time.Now().Unix(),
		FinishTime:   time.Now().Unix(),
		BasePay:      u.State.Reels.Amount,

		IsPFR:  !utils.Empty(u.State.FreespinID),
		IsDemo: u.State.IsDemo,
	}, nil
}

func BaseSessionState() SessionState {
	return SessionState{
		Username:     "",
		Balance:      0,
		Currency:     "",
		WagerID:      "some id",
		WagerAmount:  0,
		FreespinID:   "",
		SessionToken: "",
		WagerLevels:  []int64{},
		DefaultWager: 0,
		Game:         "",

		Reels: game.State{},
		//BonusStates: game.BonusStates{},

		StartBalance:  0,
		UserID:        uuid.Nil,
		Integrator:    "",
		Operator:      "",
		TransactionID: "",
		RoundID:       uuid.UUID{},
	}
}
