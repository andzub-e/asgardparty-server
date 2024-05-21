package constants

var (
	debug = "true"
	Debug = true
)

func init() {
	if debug != "true" {
		Debug = false
	}
}
