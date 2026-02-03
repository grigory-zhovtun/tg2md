package parser

import (
	"encoding/json"
	"fmt"
)

// Chat represents the root structure of Telegram export JSON.
type Chat struct {
	Name     string    `json:"name"`
	Type     string    `json:"type"`
	ID       int64     `json:"id"`
	Messages []Message `json:"messages"`
}

// Message represents a single message in the export.
type Message struct {
	ID            int64       `json:"id"`
	Type          string      `json:"type"`
	Date          string      `json:"date"`
	DateUnixtime  string      `json:"date_unixtime,omitempty"`
	From          string      `json:"from,omitempty"`
	FromID        string      `json:"from_id,omitempty"`
	ReplyToMsgID  *int64      `json:"reply_to_message_id,omitempty"`
	ForwardedFrom string      `json:"forwarded_from,omitempty"`
	Text          TextContent `json:"text"`
	TextEntities  []TextEntity `json:"text_entities,omitempty"`
	Action        string      `json:"action,omitempty"`
	Actor         string      `json:"actor,omitempty"`
}

// TextContent handles polymorphic text field (string or array of entities).
type TextContent struct {
	Plain    string
	Entities []TextEntity
}

// TextEntity represents a formatted text fragment.
// Can be either an object {"type": "...", "text": "..."} or a plain string.
type TextEntity struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Href string `json:"href,omitempty"`
}

// UnmarshalJSON handles mixed array elements (objects or plain strings).
func (te *TextEntity) UnmarshalJSON(data []byte) error {
	// Try plain string first (Telegram puts plain strings directly in array)
	var plain string
	if err := json.Unmarshal(data, &plain); err == nil {
		te.Type = "plain"
		te.Text = plain
		te.Href = ""
		return nil
	}

	// Try object
	type textEntityAlias TextEntity
	var obj textEntityAlias
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	*te = TextEntity(obj)
	return nil
}

// UnmarshalJSON handles polymorphic text field.
func (tc *TextContent) UnmarshalJSON(data []byte) error {
	// Try string first
	var plain string
	if err := json.Unmarshal(data, &plain); err == nil {
		tc.Plain = plain
		tc.Entities = nil
		return nil
	}

	// Try array of entities
	var entities []TextEntity
	if err := json.Unmarshal(data, &entities); err == nil {
		tc.Plain = ""
		tc.Entities = entities
		return nil
	}

	return fmt.Errorf("text field is neither string nor array")
}

// MarshalJSON implements json.Marshaler for TextContent.
func (tc TextContent) MarshalJSON() ([]byte, error) {
	if len(tc.Entities) > 0 {
		return json.Marshal(tc.Entities)
	}
	return json.Marshal(tc.Plain)
}

// ParseResult contains either a message or an error from parsing.
type ParseResult struct {
	Message *Message
	Error   error
}
