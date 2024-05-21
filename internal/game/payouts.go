package game

import (
	"github.com/samber/lo"
)

type Payouts struct {
	Amount int64    `json:"amount"`
	Values []Payout `json:"values"`
}

// FindPayLines - save unsorted list of payouts by pay lines at exist window.
func (p *Payouts) FindPayLines(reelWindow ReelWindow, topWindow TopWindow, figures map[FigureID]Figure) *Payouts {
	// check winLine for each symbol
	for symbol := range BasePayTable {
		chainWithCount := make([]int, 0) // count elements in each vertical line https://prnt.sc/PbUYIkI9I47S

		for x := 0; x < WindowWidth; x++ {
			if x != len(chainWithCount) {
				break
			}

			// reel window part
			for y := 4; y < WindowHeight; y++ {
				if reelWindow[x][y] != 0 { // not empty
					figureID := reelWindow[x][y]

					if figures[figureID].Symbol == symbol || figures[figureID].Symbol == "W" {
						if len(chainWithCount) != x+1 {
							chainWithCount = append(chainWithCount, 0)
						}

						chainWithCount[x]++ // was break there
					}
				}
			}

			// top window part
			topX := x - 2

			if x > 1 && x < 5 && topWindow[topX] != 0 { // if chainSize wasn't changed in reel window
				figureID := topWindow[topX]
				if figures[figureID].Symbol == symbol || figures[figureID].Symbol == "W" {
					if len(chainWithCount) != x+1 {
						chainWithCount = append(chainWithCount, 0)
					}

					chainWithCount[x]++
				}
			}
		}
		// save payout
		if len(chainWithCount) >= 4 {
			var payout Payout
			payout.Symbol = symbol
			payout.Chain = chainWithCount
			p.Values = append(p.Values, payout)
		}
	}

	// get figures ids for each payout
	for payoutID, payout := range p.Values {
		foundFiguresIDs := make(map[FigureID]bool)

		for x := 0; x < len(payout.Chain); x++ {
			// reel window
			for y := 4; y < WindowHeight; y++ {
				figureID := reelWindow[x][y]
				if figures[figureID].Symbol == payout.Symbol {
					foundFiguresIDs[figures[figureID].ID] = true
				}

				if figures[figureID].Symbol == "W" {
					foundFiguresIDs[figures[figureID].ID] = true
				}
			}
			// top window
			topX := x - 2
			if x > 1 && x < 5 && topWindow[topX] != 0 {
				figureID := topWindow[topX]
				if figures[figureID].Symbol == payout.Symbol || figures[figureID].Symbol == "W" {
					foundFiguresIDs[figures[figureID].ID] = true
				}
			}
		}
		// convert found ids from map to slice
		for figureID := range foundFiguresIDs {
			//payout.Figures = append(payout.Figures, figureID)
			p.Values[payoutID].Figures = append(p.Values[payoutID].Figures, figureID)
		}
	}

	return p
}

// DestroyPayLines destroyed pay lines
func (p *Payouts) DestroyPayLines(reelWindow *ReelWindow, topWindow *TopWindow, figures map[FigureID]Figure) *Payouts {
	for _, payout := range p.Values {
		for _, figureID := range payout.Figures {
			figure := figures[figureID]

			if figure.Y == 0 && figure.X >= 0 && figure.X < 3 {
				if topWindow[figure.X] == figureID {
					topWindow[figure.X] = 0 // delete from top window

					continue
				}
			}

			reelWindow.DeleteFigureFromWindow(figure)
		}
	}

	return p
}

// IsBonusGameTriggered - return true for 3 scatters.
func (p *Payouts) IsBonusGameTriggered(reelWindow *ReelWindow, topWindow *TopWindow, figures map[FigureID]Figure) bool {
	var (
		payout  Payout
		scatter = make([]int, WindowWidth)
	)

	// check reel window
	for x := 0; x < WindowWidth; x++ {
		for y := 4; y < WindowHeight; y++ {
			if reelWindow[x][y] != 0 { // not empty
				figureID := reelWindow[x][y]

				if figures[figureID].Symbol == "F" {
					scatter[figures[figureID].X]++

					payout.Figures = append(payout.Figures, figureID)
				}
			}
		}
	}

	// check top window
	for x := 0; x < TopWindowWidth; x++ {
		if topWindow[x] != 0 { // not empty
			figureID := topWindow[x]

			if figures[figureID].Symbol == "F" {
				scatter[figures[figureID].X]++

				payout.Figures = append(payout.Figures, figureID)
			}
		}
	}

	if lo.Sum(scatter) >= 3 {
		payout.Symbol = "F"
		payout.Chain = scatter

		p.Values = append(p.Values, payout)

		return true
	}

	return false
}

// DestroyScatters - clear scatters that trigger bonuses.
func (p *Payouts) DestroyScatters(reelWindow *ReelWindow, topWindow *TopWindow, figures map[FigureID]Figure) *Payouts {
	// for reel window
	for x := 0; x < WindowWidth; x++ {
		for y := 4; y < WindowHeight; y++ {
			if reelWindow[x][y] != 0 { // not empty
				figureID := reelWindow[x][y]

				if figures[figureID].Symbol == "F" {
					reelWindow[x][y] = 0
				}
			}
		}
	}

	// for top window
	for x := 0; x < TopWindowWidth; x++ {
		if topWindow[x] != 0 { // not empty
			figureID := topWindow[x]

			if figures[figureID].Symbol == "F" {
				topWindow[x] = 0
			}
		}
	}

	return p
}

// CalculateAmounts - get amounts by pay table
func (p *Payouts) CalculateAmounts(wager int64) *Payouts {
	for payoutID, payout := range p.Values {
		if symbolPay, ok := BasePayTable[payout.Symbol]; ok {
			var multiplayer = 1

			amount := int(symbolPay[uint(len(payout.Chain))])

			for _, count := range payout.Chain {
				multiplayer *= count
			}

			win := (wager / MathWager) * int64(amount*multiplayer)

			p.Values[payoutID].Amount = win
			p.Amount += win
		}
	}

	return p
}
