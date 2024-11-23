package utils

import (
	"regexp"
	"strings"
)

func MakeSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters (keep alphanumeric and hyphens)
	re := regexp.MustCompile("[^a-z0-9-]")
	slug = re.ReplaceAllString(slug, "")

	// Remove consecutive hyphens
	slug = regexp.MustCompile("-+").ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Trim to a maximum length of 30 characters
	if len(slug) > 30 {
		slug = slug[:30]
	}

	return slug
}
