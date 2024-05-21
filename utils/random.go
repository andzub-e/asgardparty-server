package utils

import (
	"bitbucket.org/electronicjaw/asgardparty-server/pkg/rng"
)

var rngClient rng.Client

func PatchRand(rngCli rng.Client) {
	rngClient = rngCli
}

// bounds are [left:right).
func RandInt(leftBound, rightBound int) int {
	if leftBound == rightBound {
		return 0
	}

	n, err := rngClient.Rand(uint64(rightBound - leftBound))
	if err != nil {
		// device doesn't support crypto random
		panic(err)
	}

	return int(n) + leftBound
}
