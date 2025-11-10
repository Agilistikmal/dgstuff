package pkg

import (
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func GenerateSlug(text string, withRandom bool) string {
	slug := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(text), "-")
	slug = strings.TrimSuffix(slug, "-")
	if withRandom {
		slug += "-" + strings.ToLower(uuid.New().String())[:4]
	}
	return slug
}
