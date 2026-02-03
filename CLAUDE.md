# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**tg2md** — CLI-утилита на Go для конвертации JSON-экспорта Telegram-чата в Markdown-файлы. Оптимизирована для обработки больших файлов через потоковый парсинг JSON.

## Build & Run

```bash
# Build
go build -o tg2md ./cmd/tg2md

# Run
./tg2md <input.json> [output_path]

# Run tests
go test ./...

# Run single test
go test -run TestName ./path/to/package
```

## Architecture

```
cmd/tg2md/       # Entry point, CLI argument parsing
internal/
  parser/        # JSON streaming parser for Telegram export
  converter/     # Message to Markdown conversion logic
  writer/        # File output, monthly splitting
  sanitizer/     # Text cleanup (Unicode, formatting)
  logger/        # Colored console output + error logging
```

## Key Design Decisions

- **Streaming JSON**: Use `json.Decoder` for memory-efficient parsing of large exports
- **Monthly splitting**: Output files named `{group_name}_{month}_{year}.md`
- **Cyrillic filenames**: Preserve Cyrillic in filenames, replace only spaces and special chars with `_`
- **Error handling**: Skip problematic messages silently, log to `errors.log`

## Telegram Export Format Notes

- Messages contain `text` field which can be string or array of text entities
- Text entities have `type` field: `bold`, `italic`, `code`, `text_link`, etc.
- Reply messages have `reply_to_message_id` field
- Forwarded messages have `forwarded_from` field
- Service messages have `action` field instead of `text`

## Spec Reference

Full specification in [SPEC.md](SPEC.md)
