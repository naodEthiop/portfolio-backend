package validator

import (
	"regexp"
	"strings"
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9-]`)
var multipleHyphen = regexp.MustCompile(`-+`)

func Slug(input string) string {
	slug := strings.ToLower(strings.TrimSpace(input))
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = nonSlugChars.ReplaceAllString(slug, "")
	slug = multipleHyphen.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}
