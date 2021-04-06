package neisgo

import (
	"regexp"
	"strings"
)

var (
	plainRegex = regexp.MustCompile(`[0-9]+\.`)
)

func genPlainText(text string) string {
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = plainRegex.ReplaceAllString(text, "")

	return text
}
