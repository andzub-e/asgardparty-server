package game

func GetFigureListWeight() (r int) {
	for _, figure := range FigureList {
		r += figure.Weight
	}

	return
}

var FigureList = []Figure{
	{
		Name:   "l4",
		Symbol: "A",
		Weight: 308,
		Mask: []string{
			"1",
		},
	},

	{
		Name:   "l3",
		Symbol: "B",
		Weight: 308,
		Mask: []string{
			"1",
		},
	},
	{
		Name:   "l2",
		Symbol: "C",
		Weight: 308,
		Mask: []string{
			"1",
		},
	},
	{
		Name:   "l1",
		Symbol: "D",
		Weight: 308,
		Mask: []string{
			"1",
		},
	},
	{
		Name:   "m21",
		Symbol: "O",
		Weight: 60,
		Mask: []string{
			"1",
		},
	},
	{
		Name:   "m11",
		Symbol: "P",
		Weight: 60,
		Mask: []string{
			"1",
		},
	},
	{
		Name:   "m22",
		Symbol: "O",
		Weight: 89,
		Mask: []string{
			"11",
			"11",
		},
	},
	{
		Name:   "m12",
		Symbol: "P",
		Weight: 87,
		Mask: []string{
			"11",
			"11",
		},
	},
	{
		Name:   "m23",
		Symbol: "O",
		Weight: 1,
		Mask: []string{
			"111",
			"111",
			"111",
		},
	},
	{
		Name:   "m13",
		Symbol: "P",
		Weight: 9,
		Mask: []string{
			"111",
			"111",
			"111",
		},
	},
	{
		Name:   "r1",
		Symbol: "X",
		Weight: 52,
		Mask: []string{
			"1",
			"1",
		},
	},
	{
		Name:   "r2",
		Symbol: "X",
		Weight: 60,
		Mask: []string{
			"11",
			"11",
			"11",
			"11",
		},
	},
	{
		Name:   "h",
		Symbol: "Y",
		Weight: 70,
		Mask: []string{
			"010",
			"010",
			"010",
			"111",
		},
	},
	{
		Name:   "c",
		Symbol: "Z",
		Weight: 84,
		Mask: []string{
			"010",
			"111",
			"010",
		},
	},

	// specials:

	{
		Name:   "m4",
		Symbol: "M",
		Weight: 1,
		Mask: []string{
			"1111",
			"1111",
			"1111",
			"1111",
		},
		IsSpecial: true,
	},
	{
		Name:   "m3",
		Symbol: "M",
		Weight: 2,
		Mask: []string{
			"111",
			"111",
			"111",
		},
		IsSpecial: true,
	},
	{
		Name:   "m2",
		Symbol: "M",
		Weight: 15,
		Mask: []string{
			"11",
			"11",
		},
		IsSpecial: true,
	},
	{
		Name:   "m1",
		Symbol: "M",
		Weight: 30,
		Mask: []string{
			"1",
		},
		IsSpecial: true,
	},
	{
		Name:   "f",
		Symbol: "F",
		Weight: 32,
		Mask: []string{
			"1",
		},
		IsSpecial: true,
	},
	{
		Name:   "w",
		Symbol: "W",
		Weight: 2,
		Mask: []string{
			"1",
		},
		IsSpecial: true,
	},
}
