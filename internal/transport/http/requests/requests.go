package requests

type StateRequest struct {
	Integrator string      `json:"integrator" validate:"integrator,required"`
	Game       string      `json:"game" validate:"game,required"`
	Params     interface{} `json:"params" validate:"required"`
}

type PlaceWagerRequest struct {
	SessionToken string `json:"session_token" form:"session_token" validate:"required"`
	Currency     string `json:"currency" form:"currency" validate:"required"`
	Wager        int64  `json:"wager" form:"wager" validate:"required"`
	FreeSpinID   string `json:"freespin_id"`
}

type UpdateSpinsIndexesRequest struct {
	SessionToken   string `json:"session_token" form:"session_token" validate:"required"`
	BaseSpinIndex  *int   `json:"base_spin_index" form:"base_spin_index" validate:"gte=0,required"`
	BonusSpinIndex *int   `json:"bonus_spin_index" form:"bonus_spin_index" validate:"gte=0,required"`
}

type GetFreeSpinsRequest struct {
	SessionToken string `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
}

type FreeSpinsWithIntegratorBetRequest struct {
	SessionToken    string `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
	IntegratorBetId string `json:"integrator_bet_id" form:"integrator_bet_id" query:"integrator_bet_id" validate:"required"`
}

type HistoryRequest struct {
	SessionToken string `json:"session_token" form:"session_token" query:"session_token" validate:"required"`
	Page         *int   `json:"page" form:"page" query:"page" validate:"required,gt=0"`
	Count        *int   `json:"count" form:"count" query:"count" validate:"required,gt=0"`
}

type CheatCustomFiguresRequest struct {
	SessionToken string   `json:"session_token"`
	Figures      []string `json:"figures"`
}
