package game

import (
	"fmt"
)

type TopWindow [TopWindowWidth]FigureID // ReelWindow[X][Y]

func (tw *TopWindow) DrawFigure(figure Figure) *TopWindow {
	tw[figure.X] = figure.ID

	return tw
}

func (tw *TopWindow) GetMysteryPositions(topFigures map[FigureID]Figure) (mysteryPositions []FigurePosition) {
	y := 0

	for x := 0; x < TopWindowWidth; x++ {
		figureID := tw[x]
		figure := topFigures[figureID]

		if figure.Symbol == "M" {
			var pos FigurePosition
			pos.X = x
			pos.Y = y
			mysteryPositions = append(mysteryPositions, pos)
		}
	}

	return
}

// PushAllToLeft - waterfall of figures, provide new positions for json.
func (tw *TopWindow) PushAllToLeft(figures map[FigureID]Figure, newFiguresPosition *[]FigurePosition) *TopWindow {
	var (
		tmpTopWindow [TopWindowWidth]FigureID
		offsetX      int
	)

	for x := 0; x < TopWindowWidth; x++ {
		if tw[x] != 0 {
			tmpTopWindow[offsetX] = tw[x]
			offsetX++
		}
	}

	// set new positions and save
	for x := 0; x < TopWindowWidth; x++ {
		if tw[x] != tmpTopWindow[x] {
			figureID := tmpTopWindow[x]
			if figureID != 0 {
				figure := figures[figureID]
				figure.X = x
				pos := FigurePosition{
					ID: figure.ID,
					X:  figure.X,
					Y:  figure.Y,
				}
				*newFiguresPosition = append(*newFiguresPosition, pos)
				figures[figureID] = figure
			}
		}

		tw[x] = tmpTopWindow[x]
	}

	// check for new positions
	/*
		fmt.Println("PushAllToLeft: check for new positions")
		for x := 0; x < constants.TopWindowWidth; x++ {
			figureID := tw[x]
			figure := figures[figureID]
			if figure.X != x {
				figure.X = x
				pos := FigurePosition{
					ID: figure.ID,
					X:  figure.X,
					Y:  figure.Y,
				}
				*newFiguresPosition = append(*newFiguresPosition, pos)
				figures[figureID] = figure
			}
		}
	*/

	return tw
}

func (tw *TopWindow) Print() *TopWindow {
	fmt.Printf("          ")

	for x := 0; x < len(tw); x++ {
		if tw[x] > 9 {
			fmt.Printf("[%v ]", tw[x])
		} else {
			fmt.Printf("[ %v ]", tw[x])
		}
	}
	fmt.Printf("\n")

	return tw
}
