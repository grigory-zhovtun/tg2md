package sanitizer

import (
	"testing"
)

func TestSanitizeText_RemovesZeroWidthChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "zero width space",
			input:    "Hello\u200BWorld",
			expected: "HelloWorld",
		},
		{
			name:     "zero width joiner",
			input:    "Hello\u200DWorld",
			expected: "HelloWorld",
		},
		{
			name:     "BOM",
			input:    "\uFEFFHello",
			expected: "Hello",
		},
		{
			name:     "multiple invisible chars",
			input:    "\u200B\u200CHello\u200DWorld\uFEFF",
			expected: "HelloWorld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeText(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeText_PreservesEmoji(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple emoji",
			input:    "Hello 😀 World",
			expected: "Hello 😀 World",
		},
		{
			name:     "multiple emojis",
			input:    "👍 Great job! 🎉",
			expected: "👍 Great job! 🎉",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeText(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeText_PreservesCyrillic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "russian text",
			input:    "Привет, мир!",
			expected: "Привет, мир!",
		},
		{
			name:     "mixed russian and english",
			input:    "Hello Мир World",
			expected: "Hello Мир World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeText(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName_ReplacesSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single space",
			input:    "Hello World",
			expected: "Hello_World",
		},
		{
			name:     "multiple spaces",
			input:    "Hello   World",
			expected: "Hello_World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName_ReplacesSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "forward slash",
			input:    "Hello/World",
			expected: "Hello_World",
		},
		{
			name:     "backslash",
			input:    "Hello\\World",
			expected: "Hello_World",
		},
		{
			name:     "colon",
			input:    "Hello:World",
			expected: "Hello_World",
		},
		{
			name:     "asterisk",
			input:    "Hello*World",
			expected: "Hello_World",
		},
		{
			name:     "question mark",
			input:    "Hello?World",
			expected: "Hello_World",
		},
		{
			name:     "quotes",
			input:    `Hello"World`,
			expected: "Hello_World",
		},
		{
			name:     "angle brackets",
			input:    "Hello<World>",
			expected: "Hello_World",
		},
		{
			name:     "pipe",
			input:    "Hello|World",
			expected: "Hello_World",
		},
		{
			name:     "multiple special chars",
			input:    "Hello/:*?World",
			expected: "Hello_World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName_CollapsesUnderscores(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "double underscore",
			input:    "Hello__World",
			expected: "Hello_World",
		},
		{
			name:     "triple underscore",
			input:    "Hello___World",
			expected: "Hello_World",
		},
		{
			name:     "mixed spaces and special",
			input:    "Hello / World",
			expected: "Hello_World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName_PreservesCyrillic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "russian chat name",
			input:    "Рабочий чат",
			expected: "Рабочий_чат",
		},
		{
			name:     "russian with special chars",
			input:    "Рабочий: чат",
			expected: "Рабочий_чат",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeName_TrimsLeadingTrailingUnderscores(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "leading space",
			input:    " Hello",
			expected: "Hello",
		},
		{
			name:     "trailing space",
			input:    "Hello ",
			expected: "Hello",
		},
		{
			name:     "leading special char",
			input:    "/Hello",
			expected: "Hello",
		},
		{
			name:     "trailing special char",
			input:    "Hello/",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestContainsOnlyWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: true,
		},
		{
			name:     "only tabs and newlines",
			input:    "\t\n\r",
			expected: true,
		},
		{
			name:     "only zero-width chars",
			input:    "\u200B\u200C\u200D",
			expected: true,
		},
		{
			name:     "mixed whitespace and invisible",
			input:    " \u200B \t",
			expected: true,
		},
		{
			name:     "text with whitespace",
			input:    "  Hello  ",
			expected: false,
		},
		{
			name:     "single character",
			input:    "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsOnlyWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("ContainsOnlyWhitespace(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
