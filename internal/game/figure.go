package game

import (
	"bitbucket.org/electronicjaw/asgardparty-server/utils"
)

type FigureID int64

type Figure struct {
	ID        FigureID   `json:"id"`
	Name      string     `json:"name"`
	Symbol    ReelSymbol `json:"symbol"`
	X         int        `json:"x"`
	Y         int        `json:"y"`
	Weight    int        `json:"weight"`
	Mask      []string   `json:"mask"`
	IsSpecial bool       `json:"is_special"` // true for freespins, mysteries, wilds
}

type FigurePosition struct {
	ID FigureID `json:"id"`
	X  int      `json:"x"`
	Y  int      `json:"y"`
}

func GenerateNewFigure(figureSequence *int64) (newFigure Figure) {
	weight := 0
	randomIndex := utils.RandInt(0, GetFigureListWeight())

	for _, figure := range FigureList {
		weight += figure.Weight
		if randomIndex < weight {
			newFigure = figure
			newFigure.ID = GetNextFigureID(figureSequence)

			break
		}
	}

	return
}

func GetNextFigureID(figureSequence *int64) FigureID {
	*figureSequence++

	return FigureID(*figureSequence)
}

func (f *Figure) Width() (value int) {
	value = len(f.Mask[0])

	return
}

func (f *Figure) Height() (value int) {
	value = len(f.Mask)

	return
}

func (f *Figure) IsVisible() bool {
	return f.Y+f.Height() > 4
}

func (f *Figure) IsSmall() bool {
	return f.Width() == 1 && f.Height() == 1
}

func (f *Figure) NumberOfBlocks() (value int) {
	value = len(f.Mask)

	for y := 0; y < len(f.Mask); y++ {
		line := f.Mask[y]
		for x := 0; x < len(line); x++ {
			symbol := string(line[x])
			if symbol == "1" {
				value++
			}
		}
	}

	return
}
