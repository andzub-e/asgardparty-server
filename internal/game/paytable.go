package game

type Multiplier uint
type Paytable map[ReelSymbol]PaytableLine
type PaytableLine map[uint]Multiplier
type ReelSymbol string

var BasePayTable = Paytable{
	"A": {4: 1, 5: 2, 6: 3, 7: 5},
	"B": {4: 1, 5: 2, 6: 3, 7: 5},
	"C": {4: 1, 5: 2, 6: 3, 7: 5},
	"D": {4: 1, 5: 2, 6: 3, 7: 5},
	"O": {4: 2, 5: 3, 6: 4, 7: 7},
	"P": {4: 2, 5: 3, 6: 4, 7: 7},
	"X": {4: 3, 5: 4, 6: 5, 7: 8},
	"Y": {4: 3, 5: 4, 6: 5, 7: 8},
	"Z": {4: 5, 5: 6, 6: 8, 7: 15},
}
