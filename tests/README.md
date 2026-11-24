# Тесты

Модульные тесты лежат рядом с кодом в `src/cmd/...`: Go ищет `_test.go`
только в каталоге пакета, поэтому `go test ./...` запускается из `src/`.
В этом каталоге хранятся сводка и отчёт о покрытии.

| Файл | Проверки |
| --- | --- |
| `src/cmd/words/words_test.go` | `TestGeneratePicksRequestedAmountFromList` — 25 слов из списка на 100; `TestGenerateLimitsRequestToAvailableWords` — запрос больше доступного списка |
| `src/cmd/results_io_test.go` | `TestPersistResultsAddsResultToExistingFile` — запись в существующий `results.json`; `TestPersistResultsCreatesFileForFirstResult` — создание файла для первого результата |
| `src/cmd/config_tui_test.go` | `TestToggleConfigDisablesWordList` — отключение набора в TUI сохраняется в конфигурации; `TestToggleConfigHandlesTenConsecutiveChanges` — десять переключений подряд |

## Покрытие

Обновляется командой `cd src && make test`.

<!-- coverage:start -->
Обновлено: 2026-11-20 10:41 MSK

Итоговое покрытие: **39.2%**.
Покрытие по функциям — в `coverage-func.txt`, сырой профиль — в
`coverage.out`.
<!-- coverage:end -->
