package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/grigoriizhovtun/tg2md/internal/parser"
	"github.com/grigoriizhovtun/tg2md/internal/sanitizer"
)

// Converter transforms parsed messages to Markdown format.
type Converter struct {
	messageCache map[int64]string
}

// New creates a new Converter.
func New() *Converter {
	return &Converter{
		messageCache: make(map[int64]string),
	}
}

// ConvertTextEntities converts text entities to Markdown.
func (c *Converter) ConvertTextEntities(entities []parser.TextEntity) string {
	var builder strings.Builder

	for _, entity := range entities {
		text := sanitizer.SanitizeText(entity.Text)

		switch entity.Type {
		case "bold":
			builder.WriteString("**")
			builder.WriteString(text)
			builder.WriteString("**")
		case "italic":
			builder.WriteString("_")
			builder.WriteString(text)
			builder.WriteString("_")
		case "code", "pre":
			builder.WriteString("`")
			builder.WriteString(text)
			builder.WriteString("`")
		case "text_link":
			// Use only URL, discard link text per spec
			if entity.Href != "" {
				builder.WriteString(entity.Href)
			} else {
				builder.WriteString(text)
			}
		case "link":
			// Plain URL, keep as-is
			builder.WriteString(text)
		case "mention", "hashtag", "email", "phone", "plain", "":
			// Keep text as-is
			builder.WriteString(text)
		default:
			// Unknown type, keep text
			builder.WriteString(text)
		}
	}

	return builder.String()
}

// ConvertMessage converts a parsed message to Markdown string.
// Returns the formatted line and any error.
func (c *Converter) ConvertMessage(msg *parser.Message) (string, time.Time, error) {
	// Parse timestamp
	timestamp, parsedTime, err := formatTimestamp(msg.Date)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	// Handle service messages
	if msg.Type == "service" || msg.Action != "" {
		actor := msg.Actor
		if actor == "" {
			actor = msg.From
		}
		if actor == "" {
			actor = "Unknown"
		}
		return fmt.Sprintf("[%s] [Служебное: %s %s]", timestamp, actor, msg.Action), parsedTime, nil
	}

	// Convert text content
	var text string
	if msg.Text.Plain != "" {
		text = sanitizer.SanitizeText(msg.Text.Plain)
	} else if len(msg.Text.Entities) > 0 {
		text = c.ConvertTextEntities(msg.Text.Entities)
	} else if len(msg.TextEntities) > 0 {
		text = c.ConvertTextEntities(msg.TextEntities)
	}

	// Check for empty message
	if sanitizer.ContainsOnlyWhitespace(text) {
		return "", time.Time{}, fmt.Errorf("empty message")
	}

	// Cache for reply lookups
	c.CacheMessage(msg.ID, text)

	// Build message prefix
	author := msg.From
	if author == "" {
		author = "Unknown"
	}

	var prefix string

	// Handle forwarded messages
	if msg.ForwardedFrom != "" {
		prefix = fmt.Sprintf("[Переслано от: %s] ", msg.ForwardedFrom)
	}

	// Handle replies (takes precedence over forwarded prefix)
	if msg.ReplyToMsgID != nil {
		replyText := "..."
		if cached, ok := c.GetCachedMessage(*msg.ReplyToMsgID); ok {
			replyText = truncateForReply(cached, 50)
		}
		prefix = fmt.Sprintf("[В ответ на: \"%s\"] ", replyText)
	}

	return fmt.Sprintf("[%s] %s: %s%s", timestamp, author, prefix, text), parsedTime, nil
}

// CacheMessage stores message text for reply lookups.
func (c *Converter) CacheMessage(id int64, text string) {
	c.messageCache[id] = text
}

// GetCachedMessage retrieves cached message text for replies.
func (c *Converter) GetCachedMessage(id int64) (string, bool) {
	text, ok := c.messageCache[id]
	return text, ok
}

// formatTimestamp parses ISO timestamp and formats it as [YYYY-MM-DD HH:MM].
func formatTimestamp(isoTime string) (string, time.Time, error) {
	// Try standard ISO format
	t, err := time.Parse("2006-01-02T15:04:05", isoTime)
	if err != nil {
		// Try with timezone
		t, err = time.Parse(time.RFC3339, isoTime)
		if err != nil {
			return "", time.Time{}, err
		}
	}
	return t.Format("2006-01-02 15:04"), t, nil
}

// truncateForReply truncates text for reply preview.
func truncateForReply(text string, maxLen int) string {
	// Remove newlines for cleaner preview
	text = strings.ReplaceAll(text, "\n", " ")

	if len(text) <= maxLen {
		return text
	}

	// Truncate at rune boundary
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}

	return string(runes[:maxLen-3]) + "..."
}
