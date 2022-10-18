package stringutils

import (
	"regexp"
)

const snakeCasePattern = `[a-z][a-z0-9]*(_[a-z0-9]+)*`

var snakeCaseRegex = regexp.MustCompile(`^` + snakeCasePattern + `$`)

func IsSnakeCase(s string) bool {
	return snakeCaseRegex.MatchString(s)
}
