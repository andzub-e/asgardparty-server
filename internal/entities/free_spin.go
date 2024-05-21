package entities

import (
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/overlord"
	"time"
)

type FreeSpin struct {
	ID         string    `json:"id"`
	Currency   string    `json:"currency"`
	ExpireDate time.Time `json:"expire_date"`
	Value      int       `json:"value"`
	Game       string    `json:"game"`
	SpinCount  int       `json:"spin_count"`
}

func FreeSpinsFromLord(bets []*overlord.FreeBet) []*FreeSpin {
	spins := make([]*FreeSpin, 0)

	for _, bet := range bets {
		spins = append(spins, &FreeSpin{
			ID:         bet.Id,
			Currency:   bet.Currency,
			ExpireDate: time.UnixMilli(bet.ExpireDate),
			Value:      int(bet.Value),
			Game:       bet.Game,
			SpinCount:  int(bet.SpinCount),
		})
	}

	return spins
}
