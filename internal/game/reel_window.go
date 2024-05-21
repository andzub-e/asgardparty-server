package game

import (
	"fmt"
)

type ReelWindow [WindowWidth]ReelWindowColumn // ReelWindow[X][Y]
type ReelWindowColumn [WindowHeight]FigureID

func (rw *ReelWindow) DrawFigure(figure Figure) *ReelWindow {
	for y := 0; y < figure.Height(); y++ {
		for x := 0; x < figure.Width(); x++ {
			if string(figure.Mask[y][x]) == "1" {
				rw[figure.X+x][figure.Y+y] = figure.ID
			}
		}
	}

	return rw
}

// IsExistPlaceForFigure - get place for figure and set new Figure(x,y).
func (rw *ReelWindow) PositionForFigureFound(figure *Figure, offsetX *int) bool {
	figure.X = *offsetX
	figure.Y = 4 - (figure.Height() - 1)
	tmpFigure := *figure

	for i := 1; i <= WindowWidth; i++ {
		//fmt.Printf("*offset_x: %v\n", offsetX)
		// check that figure no out of window (right side)
		if tmpFigure.X+figure.Width() > WindowWidth {
			tmpFigure.X = 0
		}

		// если позиция найдена
		if rw.PushFigureToBottom(&tmpFigure) {
			figure.X = tmpFigure.X
			figure.Y = tmpFigure.Y
			*offsetX = figure.X + figure.Width()

			return true
		}

		tmpFigure.Y = figure.Y
		tmpFigure.X++
	}

	return false
}

// PushAllToBottom - waterfall of figures, provide new positions for json.
func (rw *ReelWindow) PushAllToBottom(figures map[FigureID]Figure, newFiguresPosition *[]FigurePosition) *ReelWindow {
	//alreadyCheckedFiguresIDs := make(map[FigureID]bool)
	for y := WindowHeight - 1; y >= 0; y-- {
		for x := 0; x < WindowWidth; x++ {
			if rw[x][y] != 0 {
				figureID := rw[x][y]
				figure := figures[figureID]
				rw.DeleteFigureFromWindow(figure)

				if rw.PushFigureToBottom(&figure) {
					// save pos for changed Y
					if figures[figureID].Y != figure.Y {
						pos := FigurePosition{
							ID: figure.ID,
							X:  figure.X,
							Y:  figure.Y,
						}
						*newFiguresPosition = append(*newFiguresPosition, pos)
						figures[figureID] = figure
					}
				}

				rw.DrawFigure(figure)
			}
		}
	}

	return rw
}

// PushFigureToBottom.
func (rw *ReelWindow) PushFigureToBottom(figure *Figure) (isPushed bool) {
	tmpFigure := *figure

	for ; tmpFigure.Y < WindowHeight-(tmpFigure.Height()-1); tmpFigure.Y++ {
		if !rw.IsFreePlaceForFigure(tmpFigure) {
			break
		}

		figure.Y = tmpFigure.Y
		isPushed = true
	}

	return
}

// IsFreePlaceForFigure - check that place is free.
func (rw *ReelWindow) IsFreePlaceForFigure(figure Figure) bool {
	for y := 0; y < figure.Height(); y++ {
		for x := 0; x < figure.Width(); x++ {
			if string(figure.Mask[y][x]) == "1" {
				if rw[figure.X+x][figure.Y+y] != 0 {
					return false
				}
			}
		}
	}

	return true
}

func (rw *ReelWindow) Print(comment string) *ReelWindow {
	fmt.Println(comment)

	for y := 0; y < WindowHeight; y++ {
		for x := 0; x < WindowWidth; x++ {
			if rw[x][y] > 9 {
				fmt.Printf("[%v ]", rw[x][y])
			} else {
				fmt.Printf("[ %v ]", rw[x][y])
			}
		}
		fmt.Printf("\n")
	}

	return rw
}

func (rw *ReelWindow) DeleteFigureFromWindow(figure Figure) *ReelWindow {
	for y := 0; y < figure.Height(); y++ {
		for x := 0; x < figure.Width(); x++ {
			if string(figure.Mask[y][x]) == "1" {
				rw[figure.X+x][figure.Y+y] = 0
			}
		}
	}

	return rw
}

func (rw *ReelWindow) GetMysteryPositions(reelFigures map[FigureID]Figure) (mysteryPositions []FigurePosition) {
	for y := 0; y < WindowHeight; y++ {
		for x := 0; x < WindowWidth; x++ {
			figureID := rw[x][y]
			figure := reelFigures[figureID]

			if figure.Symbol == "M" {
				var pos FigurePosition
				pos.X = x
				pos.Y = y
				mysteryPositions = append(mysteryPositions, pos)
			}
		}
	}

	return
}
