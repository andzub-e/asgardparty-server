package simulation

type Result struct {
	Game  string `xlsx:"Game"`
	Count int    `xlsx:"Count"`
	Wager int64  `xlsx:"Wager"`
	Spent int64  `xlsx:"Spent"`

	NewLine1 string `xlsx:""`

	BonusGameCount        int `xlsx:"Bonus Game Count"`
	BonusGameWithWinCount int `xlsx:"Bonus Game With Win Count"`
	WinSpinCount          int `xlsx:"Win Spin Count"`

	NewLine2 string `xlsx:""`

	BonusGameRate float64 `xlsx:"Bonus Game Rate (Spins for One Bonus)"`
	WinSpinRate   string  `xlsx:"Win Spin Rate"`

	NewLine3 string `xlsx:""`

	CountX1   int `xlsx:"-"`
	CountX10  int `xlsx:"-"`
	CountX100 int `xlsx:"-"`

	RateX1   string `xlsx:"X1 Rate"`
	RateX10  string `xlsx:"X10 Rate"`
	RateX100 string `xlsx:"X100 Rate"`

	NewLine4 string `xlsx:""`

	BaseAward           int64 `xlsx:"Base Award"`
	BonusAward          int64 `xlsx:"Bonus Award"`
	Award               int64 `xlsx:"Award"`
	CascadeBaseAward    int64 `xlsx:"Cascade Base Award"`
	FirstStageBaseAward int64 `xlsx:"First Stage Base Award"`

	NewLine5 string `xlsx:""`

	BaseAwardSquareSum  int64 `xlsx:"Base Award Square Sum"`
	BonusAwardSquareSum int64 `xlsx:"Bonus Award Square Sum"`
	AwardSquareSum      int64 `xlsx:"Award Square Sum"`

	NewLine6 string `xlsx:""`

	BaseAwardStandardDeviation  float64 `xlsx:"Base Award Standard Deviation"`
	BonusAwardStandardDeviation float64 `xlsx:"Bonus Award Standard Deviation"`
	AwardStandardDeviation      float64 `xlsx:"Award Standard Deviation"`

	NewLine7 string `xlsx:""`

	RTP               string `xlsx:"RTP"`
	RTPBaseGame       string `xlsx:"RTP Base Game"`
	RTPBonusGame      string `xlsx:"RTP Bonus Game"`
	RTPCascadeBase    string `xlsx:"RTP Cascade Base"`
	RTPFirstStageBase string `xlsx:"RTP First Stage Base"`
}
