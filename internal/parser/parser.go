package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Parser handles streaming JSON parsing of Telegram exports.
type Parser struct {
	decoder  *json.Decoder
	file     *os.File
	chatName string
	chatType string
}

// New creates a new Parser for the given file path.
func New(filePath string) (*Parser, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	return &Parser{
		decoder: json.NewDecoder(file),
		file:    file,
	}, nil
}

// GetChatInfo extracts chat name and type from the JSON.
// Must be called before StreamMessages.
func (p *Parser) GetChatInfo() (name, chatType string, err error) {
	// Reset to beginning of file
	if _, err := p.file.Seek(0, 0); err != nil {
		return "", "", fmt.Errorf("seek file: %w", err)
	}
	p.decoder = json.NewDecoder(p.file)

	// Navigate to find name and type fields
	depth := 0
	for {
		token, err := p.decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", fmt.Errorf("read token: %w", err)
		}

		switch t := token.(type) {
		case json.Delim:
			switch t {
			case '{':
				depth++
			case '}':
				depth--
			case '[':
				// Skip arrays (like messages)
				if err := skipArray(p.decoder); err != nil {
					return "", "", err
				}
			}
		case string:
			if depth == 1 {
				switch t {
				case "name":
					var val string
					if err := p.decoder.Decode(&val); err != nil {
						return "", "", fmt.Errorf("decode name: %w", err)
					}
					p.chatName = val
				case "type":
					var val string
					if err := p.decoder.Decode(&val); err != nil {
						return "", "", fmt.Errorf("decode type: %w", err)
					}
					p.chatType = val
				}
			}
		}

		// Stop once we have both
		if p.chatName != "" && p.chatType != "" {
			break
		}
	}

	if p.chatName == "" {
		return "", "", fmt.Errorf("chat name not found in JSON")
	}

	return p.chatName, p.chatType, nil
}

// StreamMessages returns a channel that yields messages one by one.
// This enables memory-efficient processing of large files.
func (p *Parser) StreamMessages() <-chan ParseResult {
	ch := make(chan ParseResult, 100)

	go func() {
		defer close(ch)

		// Reset to beginning
		if _, err := p.file.Seek(0, 0); err != nil {
			ch <- ParseResult{Error: fmt.Errorf("seek file: %w", err)}
			return
		}
		p.decoder = json.NewDecoder(p.file)

		// Navigate to "messages" array
		foundMessages := false
		depth := 0
		for {
			token, err := p.decoder.Token()
			if err == io.EOF {
				return
			}
			if err != nil {
				ch <- ParseResult{Error: fmt.Errorf("read token: %w", err)}
				return
			}

			switch t := token.(type) {
			case json.Delim:
				switch t {
				case '{':
					depth++
				case '}':
					depth--
				case '[':
					if foundMessages {
						// This is the messages array start
						goto streamMessages
					}
					// Skip other arrays
					if err := skipArray(p.decoder); err != nil {
						ch <- ParseResult{Error: err}
						return
					}
				}
			case string:
				if depth == 1 && t == "messages" {
					foundMessages = true
				}
			}
		}

	streamMessages:
		// Stream individual messages
		for p.decoder.More() {
			var msg Message
			if err := p.decoder.Decode(&msg); err != nil {
				ch <- ParseResult{Error: fmt.Errorf("decode message: %w", err)}
				continue
			}
			ch <- ParseResult{Message: &msg}
		}
	}()

	return ch
}

// Close closes the underlying file.
func (p *Parser) Close() error {
	if p.file != nil {
		return p.file.Close()
	}
	return nil
}

// skipArray skips an entire JSON array.
func skipArray(decoder *json.Decoder) error {
	depth := 1
	for depth > 0 {
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("skip array: %w", err)
		}
		switch t := token.(type) {
		case json.Delim:
			switch t {
			case '[':
				depth++
			case ']':
				depth--
			}
		}
	}
	return nil
}
