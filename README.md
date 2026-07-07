# algo-sandbox

Hands-on implementations of the data structures interviewers keep asking about.

HTTP API on `:8084` plus a tiny CLI for offline demos.

## Why this repo exists

I wanted one place to **show** LRU, rate limiters, consistent hashing, and graph search — not just talk about them on interviews.

## CLI

```bash
go run ./cmd/algo bench-lru
go run ./cmd/algo graph-demo
```

## HTTP (sample)

```bash
make docker-up

# LRU memory backend
curl -X POST localhost:8084/api/v1/lru/put -d '{"key":"a","value":"1"}'
curl localhost:8084/api/v1/lru/get?key=a

# Top-K
curl -X POST localhost:8084/api/v1/topk -d '{"k":2,"items":[{"value":"a","score":1},{"value":"b","score":9}]}'

# Graph BFS
curl -X POST localhost:8084/api/v1/graph/path -d '{"from":"A","to":"D","graph":{"A":["B","C"],"B":["D"],"C":["D"],"D":[]}}'
```

Full complexity table: [COMPLEXITY.md](COMPLEXITY.md)

## Modules

LRU · token bucket · sliding window · consistent hash · top-K heap · BFS/DFS · worker pool

MIT
