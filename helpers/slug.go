package helpers

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	// Ubah ke lowercase
	slug := strings.ToLower(text)

	// Hapus karakter non-alfanumerik kecuali spasi
	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	// Ganti spasi dan strip ganda dengan strip tunggal
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")

	return slug
}
