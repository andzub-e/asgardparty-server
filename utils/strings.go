package utils

import "strings"

func Empty(s string) bool {
	return strings.TrimSpace(s) == ""
}
