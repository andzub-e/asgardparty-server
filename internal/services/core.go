package services

import (
	"bitbucket.org/electronicjaw/asgardparty-server/internal/entities"
	"bitbucket.org/electronicjaw/asgardparty-server/internal/game"
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
	"encoding/json"
	"errors"
	"fmt"
)

type CoreService struct{}

type LookUpOrCreateUserParams struct {
	SessionID  string
	Integrator string
	Game       string
	Params     interface{}
}

func NewCoreService() *CoreService {
	return &CoreService{}
}

func (c *CoreService) CheatCustomFigures(figures []string) (data string, err error) {
	for _, newFigureName := range figures {
		if !validateFigureByName(newFigureName) {
			err = errors.New("wrong name of figure: " + newFigureName)

			return
		}
	}

	dataJSON, err := json.Marshal(figures)
	if err != nil {
		return
	}

	return string(dataJSON), nil
}

func validateFigureByName(newFigureName string) (isValid bool) {
	for _, standardFigure := range game.FigureList {
		if standardFigure.Name == newFigureName {
			return true
		}
	}

	return false
}

func GenerateNewFigure(u *entities.User, figureSequence *int64) (newFigure game.Figure) {
	newFigure = game.GenerateNewFigure(figureSequence)

	data, ok := utils.Cache.Get(u.State.SessionToken + "-base")
	// normal generate (without cheats)
	if !ok {
		return game.GenerateNewFigure(figureSequence)
	}

	// cache generate (with cheats)
	dataString := fmt.Sprintf("%v", data)

	var names []string

	json.Unmarshal([]byte(dataString), &names)

	if len(names) > 0 {
		newFigure = getFigureByName(names[0], figureSequence)
	}

	if len(names) > 1 {
		names = names[1:]
		namesJSON, _ := json.Marshal(names)
		utils.Cache.Set(u.State.SessionToken+"-base", string(namesJSON))
	} else {
		utils.Cache.Delete(u.State.SessionToken + "-base")
	}

	return newFigure
}

func getFigureByName(name string, figureSequence *int64) (newFigure game.Figure) {
	for _, figure := range game.FigureList {
		if figure.Name == name {
			newFigure = figure
			newFigure.ID = game.GetNextFigureID(figureSequence)

			return
		}
	}

	return
}

func (c *CoreService) BaseUser() *entities.User {
	return &entities.User{
		State: entities.BaseSessionState(),
	}
}

func (c *CoreService) ExecSpin(u *entities.User) error {
	spin := c.RollBaseSpin(u)
	spin.Type = game.SpinTypeBase

	u.State.Reels.Spins = append(u.State.Reels.Spins, game.Spin{BaseSpin: spin, Type: spin.Type})
	u.State.Reels.Amount += spin.Amount

	// TODO: balance
	u.State.TotalWins = u.State.Reels.Amount
	u.State.Balance = u.State.StartBalance + u.State.Reels.Amount - u.State.WagerAmount

	return nil
}

func (c *CoreService) RollBaseSpin(u *entities.User) (baseSpin game.BaseSpin) {
	var figureSequence int64

	baseSpin.Figures = make(map[game.FigureID]game.Figure)

	// avalanche
	for stageCount := 1; stageCount < 1000; stageCount++ {
		//fmt.Println("////////////// STAGE ///////////////")
		stage := game.BaseStage{}
		stage.Multiplier = int64(stageCount)
		offsetX := 0
		isMysteryExist := false
		// add figures to Reel Window
		//fmt.Println("base: generate reel window figures")
		for {
			newFigure := GenerateNewFigure(u, &figureSequence)
			if baseSpin.ReelWindow.PositionForFigureFound(&newFigure, &offsetX) {
				if newFigure.Symbol == "M" {
					isMysteryExist = true
				}

				baseSpin.ReelWindow.DrawFigure(newFigure)
				baseSpin.Figures[newFigure.ID] = newFigure
				stage.NewReelFigures = append(stage.NewReelFigures, newFigure)

				if newFigure.IsVisible() {
					continue
				}
			}

			break
		}

		// add figures to Top Window
		//fmt.Println("base: generate top window figures")
		for x := 0; x < game.TopWindowWidth; x++ {
			if baseSpin.TopWindow[x] == 0 { // empty place
				for {
					newFigure := GenerateNewFigure(u, &figureSequence)
					newFigure.X = x
					// debug
					/*
						newFigure.Name = "w"
						newFigure.Symbol = "W"
						newFigure.Weight = 70

						newFigure.Name = "m21"
						newFigure.Symbol = "O"
						newFigure.Weight = 140
					*/
					// debug end
					if newFigure.IsSmall() {
						if newFigure.Symbol == "M" {
							isMysteryExist = true
						}

						baseSpin.TopWindow.DrawFigure(newFigure)
						baseSpin.Figures[newFigure.ID] = newFigure
						stage.NewTopFigures = append(stage.NewTopFigures, newFigure)

						break
					}
				}
			}
		}
		//baseSpin.ReelWindow.Print("1. GENERATED")
		//baseSpin.TopWindow.Print()
		// transform mystery
		//fmt.Println("base: transform mystery")
		if isMysteryExist {
			baseSpin.Stages = append(baseSpin.Stages, stage) // save prev stage without destroy because need transform mystery
			stage = game.BaseStage{}
			stage.Multiplier = int64(stageCount)
			// get 1x1 figure type for mystery
			var newFigure game.Figure

			for {
				newFigure = GenerateNewFigure(u, &figureSequence)
				if newFigure.IsSmall() && !newFigure.IsSpecial {
					break
				}
			}
			// debug
			/*
				newFigure.Name = "m11"
				newFigure.Symbol = "P"
				newFigure.Weight = 140
			*/
			// debug end

			// create figure for each mystery slot
			reelMysteryPositions := baseSpin.ReelWindow.GetMysteryPositions(baseSpin.Figures)
			for _, pos := range reelMysteryPositions {
				newFigure.ID = game.GetNextFigureID(&figureSequence)
				newFigure.X = pos.X
				newFigure.Y = pos.Y
				baseSpin.ReelWindow.DrawFigure(newFigure)
				baseSpin.Figures[newFigure.ID] = newFigure
				stage.NewReelFigures = append(stage.NewReelFigures, newFigure)
			}

			topMysteryPositions := baseSpin.TopWindow.GetMysteryPositions(baseSpin.Figures)
			for _, pos := range topMysteryPositions {
				newFigure.ID = game.GetNextFigureID(&figureSequence)
				newFigure.X = pos.X
				newFigure.Y = pos.Y
				baseSpin.TopWindow.DrawFigure(newFigure)
				baseSpin.Figures[newFigure.ID] = newFigure
				stage.NewTopFigures = append(stage.NewTopFigures, newFigure)
			}

			// fill down
			baseSpin.ReelWindow.PushAllToBottom(baseSpin.Figures, &stage.NewFiguresPosition)
		}

		//baseSpin.ReelWindow.Print("2. TRANSFORMED")
		//baseSpin.TopWindow.Print()

		// check for bonus trigger
		//fmt.Println("base: check for bonus trigger")
		if stage.Payouts.IsBonusGameTriggered(&baseSpin.ReelWindow, &baseSpin.TopWindow, baseSpin.Figures) {
			var (
				bonusGame game.State
				freeSpins = 15
			)

			stage.Payouts.DestroyScatters(&baseSpin.ReelWindow, &baseSpin.TopWindow, baseSpin.Figures)

			for freeSpins > 0 {
				freeSpins--

				bonusSpin := c.RollBonusSpin(u)
				bonusSpin.FreeSpinsLeft = freeSpins
				bonusGame.Amount += bonusSpin.Amount
				freeSpins += bonusSpin.FreeSpins
				bonusSpin.Type = game.SpinTypeBonus

				bonusGame.Spins = append(bonusGame.Spins, game.Spin{BonusSpin: bonusSpin, Type: bonusSpin.Type})
			}

			stage.BonusGame = bonusGame
			baseSpin.Amount += bonusGame.Amount // save bonus game reward
		}

		// get paylines and destroy them
		//fmt.Println("base: FindPayLines")
		stage.Payouts.FindPayLines(baseSpin.ReelWindow, baseSpin.TopWindow, baseSpin.Figures)
		//fmt.Println("base: DestroyPayLines")
		stage.Payouts.DestroyPayLines(&baseSpin.ReelWindow, &baseSpin.TopWindow, baseSpin.Figures)
		//fmt.Println("base: ReelWindow PushAllToBottom")
		baseSpin.ReelWindow.PushAllToBottom(baseSpin.Figures, &stage.NewFiguresPosition)
		//fmt.Println("base: TopWindow PushAllToLeft")
		baseSpin.TopWindow.PushAllToLeft(baseSpin.Figures, &stage.NewFiguresPosition)
		// calculate amounts
		//fmt.Println("base: CalculateAmounts")
		stage.Payouts.CalculateAmounts(u.State.WagerAmount)
		baseSpin.Amount += stage.Payouts.Amount

		baseSpin.Stages = append(baseSpin.Stages, stage)

		if len(stage.Payouts.Values) == 0 {
			break
		}
	}

	return baseSpin
}

func (c *CoreService) RollBonusSpin(u *entities.User) (bonusSpin game.BonusSpin) {
	var figureSequence int64

	bonusSpin.Figures = make(map[game.FigureID]game.Figure)

	// avalanche
	for stageCount := 1; stageCount < 1000; stageCount++ {
		//fmt.Println("////////////// STAGE: BONUS ///////////////")
		stage := game.BonusStage{}
		stage.Multiplier = int64(stageCount * 3)
		offsetX := 0
		isMysteryExist := false
		// add figures to Reel Window
		//fmt.Println("bonus: generate reel window figures")
		for {
			newFigure := GenerateNewFigure(u, &figureSequence)
			if bonusSpin.ReelWindow.PositionForFigureFound(&newFigure, &offsetX) {
				if newFigure.Symbol == "M" {
					isMysteryExist = true
				}

				bonusSpin.ReelWindow.DrawFigure(newFigure)
				bonusSpin.Figures[newFigure.ID] = newFigure
				stage.NewReelFigures = append(stage.NewReelFigures, newFigure)

				if newFigure.IsVisible() {
					continue
				}
			}

			break
		}

		// add figures to Top Window
		//fmt.Println("bonus: generate top window figures")
		for x := 0; x < game.TopWindowWidth; x++ {
			if bonusSpin.TopWindow[x] == 0 { // empty place
				for {
					newFigure := GenerateNewFigure(u, &figureSequence)
					newFigure.X = x

					if newFigure.IsSmall() {
						if newFigure.Symbol == "M" {
							isMysteryExist = true
						}

						bonusSpin.TopWindow.DrawFigure(newFigure)
						bonusSpin.Figures[newFigure.ID] = newFigure
						stage.NewTopFigures = append(stage.NewTopFigures, newFigure)

						break
					}
				}
			}
		}

		// transform mystery
		//fmt.Println("bonus: transform mystery")
		if isMysteryExist {
			bonusSpin.Stages = append(bonusSpin.Stages, stage) // save prev stage without destroy because need transform mystery
			stage = game.BonusStage{}
			stage.Multiplier = int64(stageCount * 3)
			// get 1x1 figure type for mystery
			var newFigure game.Figure

			for {
				newFigure = GenerateNewFigure(u, &figureSequence)
				if newFigure.IsSmall() && !newFigure.IsSpecial {
					break
				}
			}

			mysteryPositions := bonusSpin.ReelWindow.GetMysteryPositions(bonusSpin.Figures)
			// create figure for each mystery slot
			for _, pos := range mysteryPositions {
				newFigure.ID = game.GetNextFigureID(&figureSequence)
				newFigure.X = pos.X
				newFigure.Y = pos.Y
				bonusSpin.ReelWindow.DrawFigure(newFigure)
				bonusSpin.Figures[newFigure.ID] = newFigure
				stage.NewReelFigures = append(stage.NewReelFigures, newFigure)
			}

			topMysteryPositions := bonusSpin.TopWindow.GetMysteryPositions(bonusSpin.Figures)
			for _, pos := range topMysteryPositions {
				newFigure.ID = game.GetNextFigureID(&figureSequence)
				newFigure.X = pos.X
				newFigure.Y = pos.Y
				bonusSpin.TopWindow.DrawFigure(newFigure)
				bonusSpin.Figures[newFigure.ID] = newFigure
				stage.NewTopFigures = append(stage.NewTopFigures, newFigure)
			}

			// fill down
			bonusSpin.ReelWindow.PushAllToBottom(bonusSpin.Figures, &stage.NewFiguresPosition)
		}

		// additional bonus spins trigger
		//fmt.Println("bonus: IsBonusGameTriggered")
		if stage.Payouts.IsBonusGameTriggered(&bonusSpin.ReelWindow, &bonusSpin.TopWindow, bonusSpin.Figures) {
			stage.Payouts.DestroyScatters(&bonusSpin.ReelWindow, &bonusSpin.TopWindow, bonusSpin.Figures)
			bonusSpin.FreeSpins += 15
		}

		// get paylines and destroy them
		//fmt.Println("bonus: FindPayLines")
		stage.Payouts.FindPayLines(bonusSpin.ReelWindow, bonusSpin.TopWindow, bonusSpin.Figures)
		//fmt.Println("bonus: DestroyPayLines")
		stage.Payouts.DestroyPayLines(&bonusSpin.ReelWindow, &bonusSpin.TopWindow, bonusSpin.Figures)
		//fmt.Println("bonus: ReelWindow PushAllToBottom")
		bonusSpin.ReelWindow.PushAllToBottom(bonusSpin.Figures, &stage.NewFiguresPosition)
		//fmt.Println("bonus: TopWindow PushAllToLeft")
		bonusSpin.TopWindow.PushAllToLeft(bonusSpin.Figures, &stage.NewFiguresPosition)

		// calculate amounts
		//fmt.Println("bonus: CalculateAmounts")
		stage.Payouts.CalculateAmounts(u.State.WagerAmount)
		bonusSpin.Amount += stage.Payouts.Amount

		bonusSpin.Stages = append(bonusSpin.Stages, stage)

		if len(stage.Payouts.Values) == 0 {
			break
		}
	}

	return bonusSpin
}
