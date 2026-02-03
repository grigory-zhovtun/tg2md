# Роль
Senior Go разработчик с опытом создания CLI-утилит, потоковой обработки JSON и работы с Git/GitHub

# Задача
Реализовать CLI-утилиту `tg2md` на Go для конвертации JSON-экспорта Telegram-чата в Markdown-файлы с разбивкой по месяцам. После завершения разработки создать публичный GitHub-репозиторий и запушить код.

# Подзадачи

## Инициализация проекта
1. Создать структуру директорий согласно архитектуре:
   - `cmd/tg2md/` — точка входа
   - `internal/parser/` — потоковый JSON-парсер
   - `internal/converter/` — конвертация в Markdown
   - `internal/writer/` — запись файлов с разбивкой по месяцам
   - `internal/sanitizer/` — очистка текста
   - `internal/logger/` — цветной вывод в консоль
2. Инициализировать Go-модуль (`go mod init`)
3. Создать CLAUDE.md и SPEC.md

## Реализация модулей

### Logger (`internal/logger/`)
4. Реализовать цветной консольный вывод:
   - Зелёный — успех, статистика
   - Жёлтый — предупреждения
   - Красный — ошибки
   - Синий — информация о прогрессе
5. Реализовать запись в `errors.log`

### Sanitizer (`internal/sanitizer/`)
6. Удаление невидимых Unicode-символов (zero-width и т.п.)
7. Санитизация имени группы: пробелы и спецсимволы → `_`, схлопывание множественных `_`

### Parser (`internal/parser/`)
8. Реализовать потоковый парсинг JSON через `json.Decoder`
9. Парсинг поля `text` (строка или массив text entities)
10. Обработка типов entities: `bold`, `italic`, `code`, `text_link`
11. Обработка `reply_to_message_id`, `forwarded_from`, `action`

### Converter (`internal/converter/`)
12. Конвертация text entities в Markdown-форматирование:
    - Жирный → `**текст**`
    - Курсив → `_текст_`
    - Код → `` `код` ``
    - Ссылки → только URL
13. Форматирование сообщений:
    - Обычное: `[YYYY-MM-DD HH:MM] Автор: Текст`
    - Ответ: `[В ответ на: "..."]`
    - Пересланное: `[Переслано от: ...]`
    - Служебное: `[Служебное: ...]`

### Writer (`internal/writer/`)
14. Создание директории с именем группы
15. Разбивка по месяцам: `название_группы_month_year.md`
16. Запись сообщений в соответствующие файлы

### CLI (`cmd/tg2md/`)
17. Парсинг аргументов: `<input.json> [output_path]`
18. Валидация входного файла
19. Оркестрация всех модулей
20. Вывод итоговой статистики

## Тестирование
21. Написать unit-тесты для каждого модуля
22. Проверить на тестовых JSON-файлах

## Git и публикация
23. Инициализировать локальный git-репозиторий
24. Создать публичный репозиторий на GitHub через `gh repo create`
25. Запушить код в удалённый репозиторий

# Релевантный контекст

## Архитектура
```
cmd/tg2md/       # Entry point, CLI argument parsing
internal/
  parser/        # JSON streaming parser for Telegram export
  converter/     # Message to Markdown conversion logic
  writer/        # File output, monthly splitting
  sanitizer/     # Text cleanup (Unicode, formatting)
  logger/        # Colored console output + error logging
```

## Команды сборки и запуска
```bash
go build -o tg2md ./cmd/tg2md
./tg2md <input.json> [output_path]
go test ./...
```

## Формат Telegram-экспорта
- Поле `text` — строка или массив text entities
- Text entities: `type` = `bold`, `italic`, `code`, `text_link`
- Ответы: `reply_to_message_id`
- Пересылки: `forwarded_from`
- Служебные: `action` вместо `text`

# Ограничения

- **Язык**: Go 1.21+
- **Парсинг JSON**: только потоковый через `json.Decoder` (для больших файлов)
- **Зависимости**: минимум внешних библиотек
- **Кроссплатформенность**: Linux, macOS, Windows
- **Кириллица**: сохранять в именах файлов
- **Обработка ошибок**: пропускать проблемные сообщения, логировать в `errors.log`
- **Критические ошибки**: файл не найден или невалидный JSON → завершение работы