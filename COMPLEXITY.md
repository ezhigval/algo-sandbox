# Complexity notes

**English** · [Русский](COMPLEXITY.ru.md)

Reference for interview prep — implementations in this repo.

| Module | Operation | Time | Space |
|---|---|---|---|
| LRU (memory) | Get/Put | O(1) | O(n) |
| LRU (redis) | Get/Put | O(log n) zset ops | O(n) |
| Token bucket | Allow | O(1) | O(1) |
| Sliding window | Allow | O(k) k=events in window | O(k) |
| Consistent hash | Lookup | O(log v) v=virtual nodes | O(v) |
| Top-K (heap) | Build | O(n log k) | O(k) |
| Graph BFS | Path | O(V+E) | O(V) |
| Worker pool | Submit | O(1) queue | O(w) workers |

Run benchmarks:

```bash
go test ./internal/algo/lru -bench=. -benchmem
go test ./internal/algo/ratelimit -bench=. -benchmem
```
