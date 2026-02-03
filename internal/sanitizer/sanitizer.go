package sanitizer

import (
	"regexp"
	"strings"
	"unicode"
)

// invisibleChars contains zero-width and other invisible Unicode characters to remove.
var invisibleChars = map[rune]bool{
	'\u200B': true, // Zero Width Space
	'\u200C': true, // Zero Width Non-Joiner
	'\u200D': true, // Zero Width Joiner
	'\uFEFF': true, // Zero Width No-Break Space (BOM)
	'\u2060': true, // Word Joiner
	'\u2061': true, // Function Application
	'\u2062': true, // Invisible Times
	'\u2063': true, // Invisible Separator
	'\u2064': true, // Invisible Plus
	'\u180E': true, // Mongolian Vowel Separator
}

// specialCharsPattern matches characters invalid in filenames.
var specialCharsPattern = regexp.MustCompile(`[/\\:*?"<>|]`)

// multipleUnderscoresPattern matches consecutive underscores.
var multipleUnderscoresPattern = regexp.MustCompile(`_+`)

// SanitizeText removes invisible Unicode characters from text.
// Preserves emojis, Cyrillic, and other visible characters.
func SanitizeText(text string) string {
	var builder strings.Builder
	builder.Grow(len(text))

	for _, r := range text {
		if !invisibleChars[r] {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// SanitizeName converts a group name to a valid filename.
// - Spaces → underscore
// - Special chars (/\:*?"<>|) → underscore
// - Cyrillic preserved
// - Multiple underscores collapsed to one
// - Leading/trailing underscores removed
func SanitizeName(name string) string {
	// Replace spaces with underscores
	result := strings.ReplaceAll(name, " ", "_")

	// Replace special characters with underscores
	result = specialCharsPattern.ReplaceAllString(result, "_")

	// Collapse multiple underscores
	result = multipleUnderscoresPattern.ReplaceAllString(result, "_")

	// Trim leading/trailing underscores
	result = strings.Trim(result, "_")

	return result
}

// ContainsOnlyWhitespace checks if text contains only whitespace or invisible characters.
func ContainsOnlyWhitespace(text string) bool {
	for _, r := range text {
		if !unicode.IsSpace(r) && !invisibleChars[r] {
			return false
		}
	}
	return true
}
