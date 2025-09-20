package utils

import "strings"

func Slug(phrase string, alt ...string) (slug string) {
	replaceWith := "-"

	if len(alt) > 0 {
		replaceWith = alt[0]
	}

	slug = strings.ToLower(phrase)
	slug = strings.ReplaceAll(slug, " ", replaceWith)
	return
}
