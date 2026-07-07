# Заметки по сложности

**[English](COMPLEXITY.md)** · Русский

Справочник для подготовки к собесам — реализации в этом репозитории.

| Модуль | Операция | Время | Память |
|---|---|---|---|
| LRU (memory) | Get/Put | O(1) | O(n) |
| LRU (redis) | Get/Put | O(log n) zset ops | O(n) |
| Token bucket | Allow | O(1) | O(1) |
| Sliding window | Allow | O(k) k=события в окне | O(k) |
| Consistent hash | Lookup | O(log v) v=виртуальные ноды | O(v) |
| Top-K (heap) | Build | O(n log k) | O(k) |
| Graph BFS | Path | O(V+E) | O(V) |
| Worker pool | Submit | O(1) queue | O(w) workers |

Запуск бенчмарков:

```bash
go test ./internal/algo/lru -bench=. -benchmem
go test ./internal/algo/ratelimit -bench=. -benchmem
```
