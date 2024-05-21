package entities

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/game"
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/history"
	"encoding/json"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type History struct {
	ID             uuid.UUID `json:"round_id" gorm:"primaryKey;column:id"`
	Game           string    `json:"game"`
	GameID         uuid.UUID `json:"game_id"`
	Integrator     string    `json:"integrator"`
	Operator       string    `json:"operator"`
	Provider       string    `json:"provider"`
	UserID         uuid.UUID `json:"user_id"`
	ExternalUserID string    `json:"external_user_id"`
	SessionToken   uuid.UUID `json:"session_token"`
	TransactionID  uuid.UUID `json:"transaction_id"`
	Currency       string    `json:"currency"`

	StartBalance int64 `json:"balance" gorm:"column:start_balance"`
	EndBalance   int64 `json:"final_balance" gorm:"column:end_balance"`
	Bet          int64 `json:"bet" gorm:"column:wager"`
	BasePay      int64 `json:"base_pay" gorm:"column:base_award"`
	BonusPay     int64 `json:"bonus_pay" gorm:"column:bonus_award"`

	ReelsState game.State `json:"reels" gorm:"serializer:json;column:spin"`
	StartTime  int64      `json:"start_time" gorm:"serializer:unixtime;column:created_at"`
	FinishTime int64      `json:"finish_time" gorm:"serializer:unixtime;column:updated_at"`

	SpinIndexes SpinIndexes `json:"restoring_indexes" gorm:"serializer:json;column:restoring_indexes"`

	IsShown bool `json:"is_shown"`
	IsPFR   bool `json:"is_pfr"`
	IsDemo  bool `json:"is_demo"`
}

type SpinIndexes struct {
	BaseSpinIndex  int `json:"base_spin_index"`
	BonusSpinIndex int `json:"bonus_spin_index"`
}

type HistoryPagination struct {
	Spins       []*History `json:"spins_history"`
	CurrentPage int        `json:"current_page"`
	Count       int        `json:"count"`
	Total       int        `json:"total"`
}

func (h *History) ToHistoryIn(metaData *PlayerMetaData) (*history.SpinIn, error) {
	restoring, err := json.Marshal(h.SpinIndexes)
	if err != nil {
		return nil, err
	}

	details, err := json.Marshal(h.ReelsState)
	if err != nil {
		return nil, err
	}

	return &history.SpinIn{
		CreatedAt: timestamppb.New(time.Unix(h.StartTime, 0)),
		UpdatedAt: timestamppb.New(time.Unix(h.FinishTime, 0)),

		Host:      metaData.Host,
		ClientIp:  metaData.IP,
		UserAgent: metaData.UserAgent,
		Request:   metaData.Request,

		Id:             h.ID.String(),
		GameId:         h.GameID.String(),
		Game:           h.Game,
		SessionToken:   h.SessionToken.String(),
		TransactionId:  h.TransactionID.String(),
		Integrator:     h.Integrator,
		Operator:       h.Operator,
		Provider:       h.Provider,
		InternalUserId: h.UserID.String(),
		ExternalUserId: h.ExternalUserID,

		Currency: h.Currency,

		StartBalance: uint64(h.StartBalance),
		EndBalance:   uint64(h.EndBalance),
		Wager:        uint64(h.Bet),
		BaseAward:    uint64(h.BasePay),
		BonusAward:   uint64(h.BonusPay),
		FinalAward:   uint64(h.BasePay + h.BonusPay),

		RestoringIndexes: restoring,
		Details:          details,

		IsPfr:   h.IsPFR,
		IsShown: h.IsShown,
		IsDemo:  &h.IsDemo,
	}, nil
}

func (History) TableName() string {
	return "history_records"
}

func SpinOutToHistory(out *history.SpinOut) (*History, error) {
	id, err := uuid.Parse(out.Id)
	if err != nil {
		return nil, err
	}

	internalUserID, err := uuid.Parse(out.InternalUserId)
	if err != nil {
		return nil, err
	}

	sessionToken, err := uuid.Parse(out.SessionToken)
	if err != nil {
		return nil, err
	}

	transactionID, err := uuid.Parse(out.TransactionId)
	if err != nil {
		return nil, err
	}

	gameID, err := uuid.Parse(out.GameId)
	if err != nil {
		return nil, err
	}

	gs := game.State{}
	if err := json.Unmarshal(out.Details, &gs); err != nil {
		return nil, err
	}

	si := SpinIndexes{}
	if err := json.Unmarshal(out.RestoringIndexes, &si); err != nil {
		return nil, err
	}

	return &History{
		ID:             id,
		Game:           out.Game,
		GameID:         gameID,
		Integrator:     out.Integrator,
		Operator:       out.Operator,
		Provider:       out.Provider,
		UserID:         internalUserID,
		ExternalUserID: out.ExternalUserId,
		SessionToken:   sessionToken,
		TransactionID:  transactionID,
		Currency:       out.Currency,

		StartBalance: int64(out.StartBalance),
		EndBalance:   int64(out.EndBalance),
		Bet:          int64(out.Wager),
		BasePay:      int64(out.BaseAward),
		BonusPay:     int64(out.BonusAward),
		ReelsState:   gs,
		StartTime:    out.CreatedAt.AsTime().Unix(),
		FinishTime:   out.UpdatedAt.AsTime().Unix(),

		SpinIndexes: si,

		IsShown: *out.IsShown,
		IsPFR:   *out.IsPfr,
		IsDemo:  *out.IsDemo,
	}, nil
}
