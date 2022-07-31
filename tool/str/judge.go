package str

import "strings"

func Empty(str string) bool {
	return len(str) == 0
}

func EmptyTrim(str string) bool {
	return Empty(strings.TrimSpace(str))
}
