package converter

import (
	"strings"
	"testing"

	"github.com/grigoriizhovtun/tg2md/internal/parser"
)

func TestConvertTextEntities_Bold(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "bold", Text: "важный"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "**важный**"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_Italic(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "italic", Text: "курсив"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "_курсив_"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_Code(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "code", Text: "fmt.Println()"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "`fmt.Println()`"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_Pre(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "pre", Text: "code block"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "`code block`"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_TextLink(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "text_link", Text: "ссылка", Href: "https://example.com"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "https://example.com"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_TextLinkWithoutHref(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "text_link", Text: "ссылка"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "ссылка"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_Link(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "link", Text: "https://example.com"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "https://example.com"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_Mixed(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "plain", Text: "Это "},
		{Type: "bold", Text: "важный"},
		{Type: "plain", Text: " текст с "},
		{Type: "text_link", Text: "ссылкой", Href: "https://example.com"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "Это **важный** текст с https://example.com"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertTextEntities_UnknownType(t *testing.T) {
	c := New()
	entities := []parser.TextEntity{
		{Type: "unknown_type", Text: "сохраняется"},
	}

	result := c.ConvertTextEntities(entities)
	expected := "сохраняется"

	if result != expected {
		t.Errorf("ConvertTextEntities() = %q, want %q", result, expected)
	}
}

func TestConvertMessage_Regular(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:   1,
		Type: "message",
		Date: "2024-01-15T14:30:00",
		From: "Иван",
		Text: parser.TextContent{Plain: "Привет, как дела?"},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	expected := "[2024-01-15 14:30] Иван: Привет, как дела?"

	if result != expected {
		t.Errorf("ConvertMessage() = %q, want %q", result, expected)
	}
}

func TestConvertMessage_Reply(t *testing.T) {
	c := New()

	// First, add original message to cache
	c.CacheMessage(1, "Привет, как дела?")

	replyToID := int64(1)
	msg := &parser.Message{
		ID:           2,
		Type:         "message",
		Date:         "2024-01-15T14:31:00",
		From:         "Мария",
		ReplyToMsgID: &replyToID,
		Text:         parser.TextContent{Plain: "Отлично!"},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	expected := `[2024-01-15 14:31] Мария: [В ответ на: "Привет, как дела?"] Отлично!`

	if result != expected {
		t.Errorf("ConvertMessage() = %q, want %q", result, expected)
	}
}

func TestConvertMessage_ReplyUncached(t *testing.T) {
	c := New()

	replyToID := int64(999)
	msg := &parser.Message{
		ID:           2,
		Type:         "message",
		Date:         "2024-01-15T14:31:00",
		From:         "Мария",
		ReplyToMsgID: &replyToID,
		Text:         parser.TextContent{Plain: "Отлично!"},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	// Should show "..." for uncached reply
	if !strings.Contains(result, `[В ответ на: "..."]`) {
		t.Errorf("ConvertMessage() = %q, should contain uncached reply marker", result)
	}
}

func TestConvertMessage_Forwarded(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:            3,
		Type:          "message",
		Date:          "2024-01-15T14:32:00",
		From:          "Пётр",
		ForwardedFrom: "Алексей",
		Text:          parser.TextContent{Plain: "Важная информация"},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	expected := "[2024-01-15 14:32] Пётр: [Переслано от: Алексей] Важная информация"

	if result != expected {
		t.Errorf("ConvertMessage() = %q, want %q", result, expected)
	}
}

func TestConvertMessage_Service(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:     5,
		Type:   "service",
		Date:   "2024-01-15T14:34:00",
		Actor:  "Иван",
		Action: "invite_members",
		Text:   parser.TextContent{Plain: ""},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	expected := "[2024-01-15 14:34] [Служебное: Иван invite_members]"

	if result != expected {
		t.Errorf("ConvertMessage() = %q, want %q", result, expected)
	}
}

func TestConvertMessage_ServiceWithFromField(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:     5,
		Type:   "service",
		Date:   "2024-01-15T14:34:00",
		From:   "Мария",
		Action: "join_group",
		Text:   parser.TextContent{Plain: ""},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	// Should use From field when Actor is empty
	if !strings.Contains(result, "Мария") {
		t.Errorf("ConvertMessage() = %q, should contain From field", result)
	}
}

func TestConvertMessage_EmptyText_ReturnsError(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:   1,
		Type: "message",
		Date: "2024-01-15T14:30:00",
		From: "Иван",
		Text: parser.TextContent{Plain: ""},
	}

	_, _, err := c.ConvertMessage(msg)
	if err == nil {
		t.Error("Expected error for empty message")
	}
}

func TestConvertMessage_WhitespaceOnlyText_ReturnsError(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:   1,
		Type: "message",
		Date: "2024-01-15T14:30:00",
		From: "Иван",
		Text: parser.TextContent{Plain: "   \t\n  "},
	}

	_, _, err := c.ConvertMessage(msg)
	if err == nil {
		t.Error("Expected error for whitespace-only message")
	}
}

func TestConvertMessage_MissingAuthor_UsesUnknown(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:   1,
		Type: "message",
		Date: "2024-01-15T14:30:00",
		From: "",
		Text: parser.TextContent{Plain: "Сообщение без автора"},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	if !strings.Contains(result, "Unknown:") {
		t.Errorf("ConvertMessage() = %q, should use 'Unknown' as author", result)
	}
}

func TestConvertMessage_InvalidDate_ReturnsError(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:   1,
		Type: "message",
		Date: "invalid-date",
		From: "Иван",
		Text: parser.TextContent{Plain: "Привет"},
	}

	_, _, err := c.ConvertMessage(msg)
	if err == nil {
		t.Error("Expected error for invalid date")
	}
}

func TestConvertMessage_WithTextEntities(t *testing.T) {
	c := New()
	msg := &parser.Message{
		ID:   1,
		Type: "message",
		Date: "2024-01-15T14:30:00",
		From: "Анна",
		Text: parser.TextContent{
			Entities: []parser.TextEntity{
				{Type: "plain", Text: "Это "},
				{Type: "bold", Text: "важный"},
				{Type: "plain", Text: " текст"},
			},
		},
	}

	result, _, err := c.ConvertMessage(msg)
	if err != nil {
		t.Fatalf("ConvertMessage failed: %v", err)
	}

	expected := "[2024-01-15 14:30] Анна: Это **важный** текст"

	if result != expected {
		t.Errorf("ConvertMessage() = %q, want %q", result, expected)
	}
}

func TestTruncateForReply(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "short text",
			text:     "Hello",
			maxLen:   50,
			expected: "Hello",
		},
		{
			name:     "exact length",
			text:     "Hello World",
			maxLen:   11,
			expected: "Hello World",
		},
		{
			name:     "long text",
			text:     "This is a very long message that should be truncated",
			maxLen:   20,
			expected: "This is a very lo...",
		},
		{
			name:     "with newlines",
			text:     "Hello\nWorld",
			maxLen:   50,
			expected: "Hello World",
		},
		{
			name:     "cyrillic text",
			text:     "Привет, как дела?",
			maxLen:   10,
			expected: "Привет,...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateForReply(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateForReply(%q, %d) = %q, want %q", tt.text, tt.maxLen, result, tt.expected)
			}
		})
	}
}

func TestCacheMessage(t *testing.T) {
	c := New()

	c.CacheMessage(123, "Test message")

	text, ok := c.GetCachedMessage(123)
	if !ok {
		t.Error("Expected to find cached message")
	}
	if text != "Test message" {
		t.Errorf("Cached text = %q, want %q", text, "Test message")
	}

	_, ok = c.GetCachedMessage(456)
	if ok {
		t.Error("Expected not to find non-existent message")
	}
}
