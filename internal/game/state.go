package game

import (
	"encoding/json"
	"fmt"
)

type State struct {
	Spins        []Spin `json:"spins"` // BaseSpin{} or BonusSpin{}
	IsCheatStops bool   `json:"is_cheat_stops"`
	IsAutospin   bool   `json:"is_autospin"`
	IsTurbospin  bool   `json:"is_turbospin"`
	Amount       int64  `json:"amount"`
}

type SpinType string

const (
	SpinTypeBonus SpinType = "bonus"
	SpinTypeBase  SpinType = "base"
)

type Spin struct {
	BaseSpin  BaseSpin
	BonusSpin BonusSpin
	Type      SpinType
}

func (s *Spin) UnmarshalJSON(data []byte) error {
	var tp map[string]any
	if err := json.Unmarshal(data, &tp); err != nil {
		return err
	}

	g, _ := tp["type"].(string)
	s.Type = SpinType(g)

	switch s.Type {
	case SpinTypeBonus:
		return json.Unmarshal(data, &s.BonusSpin)
	case SpinTypeBase:
		return json.Unmarshal(data, &s.BaseSpin)
	default:
		return fmt.Errorf(`unknown type: "%s"`, s.Type)
	}
}

func (s *Spin) MarshalJSON() ([]byte, error) {
	switch s.Type {
	case SpinTypeBonus:
		return json.Marshal(s.BonusSpin)
	case SpinTypeBase:
		return json.Marshal(s.BaseSpin)
	default:
		return nil, fmt.Errorf(`unknown type: "%s"`, s.Type)
	}
}

type BaseSpin struct {
	ReelWindow ReelWindow          `json:"-"`
	TopWindow  TopWindow           `json:"-"`
	Figures    map[FigureID]Figure `json:"-"`

	Stages []BaseStage `json:"stages"`
	Amount int64       `json:"amount" swaggertype:"integer" example:"100"`

	Type SpinType `json:"type"`
}

func (s *BaseSpin) BaseAward() (award int64) {
	for _, stage := range s.Stages {
		award += stage.Payouts.Amount
	}

	return award
}

func (s *BaseSpin) BonusAward() (award int64) {
	for _, stage := range s.Stages {
		if stage.BonusGame != nil {
			state := stage.BonusGame.(State)
			award += state.Amount
		}
	}

	return award
}

func (s *BaseSpin) CascadeBaseAward() (award int64) {
	for i, stage := range s.Stages {
		if i > 0 {
			award += stage.Payouts.Amount
		}
	}

	return award
}

func (s *BaseSpin) FirstStageBaseAward() (award int64) {
	for i, stage := range s.Stages {
		if i == 0 {
			award += stage.Payouts.Amount
		}
	}

	return award
}

func (s *BaseSpin) BonusGameCount() (count int) {
	for _, stage := range s.Stages {
		if stage.BonusGame != nil {
			count++
		}
	}

	return count
}

func (s *BaseSpin) BonusGameWithWinCount() (count int) {
	for _, stage := range s.Stages {
		if stage.BonusGame != nil {
			state := stage.BonusGame.(State)
			if state.Amount > 0 {
				count++
			}
		}
	}

	return count
}

type BaseStage struct {
	Multiplier         int64            `json:"multiplier"`
	NewReelFigures     []Figure         `json:"new_reel_figures"`
	NewTopFigures      []Figure         `json:"new_top_figures"`
	NewFiguresPosition []FigurePosition `json:"new_figures_position"`
	Payouts            Payouts          `json:"payouts"`
	BonusGame          interface{}      `json:"bonus_game"` // State{}
}

type BonusSpin struct {
	ReelWindow ReelWindow          `json:"-"`
	TopWindow  TopWindow           `json:"-"`
	Figures    map[FigureID]Figure `json:"-"`

	Stages        []BonusStage `json:"stages"`
	Amount        int64        `json:"amount" swaggertype:"integer" example:"100"`
	FreeSpinsLeft int          `json:"free_spins_left"`
	FreeSpins     int          `json:"-"`

	Type SpinType `json:"type"`
}

type BonusStage struct {
	Multiplier         int64            `json:"multiplier"`
	NewReelFigures     []Figure         `json:"new_reel_figures"`
	NewTopFigures      []Figure         `json:"new_top_figures"`
	NewFiguresPosition []FigurePosition `json:"new_figures_position"`
	Payouts            Payouts          `json:"payouts"`
}

type Payout struct {
	Symbol  ReelSymbol `json:"symbol" example:"A"`
	Chain   []int      `json:"count" example:"5"`
	Figures []FigureID `json:"figures" example:"5"`
	Amount  int64      `json:"amount" swaggertype:"integer" example:"100"`
}
