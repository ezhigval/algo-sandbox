# algo-sandbox

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
[![CI](https://github.com/ezhigval/algo-sandbox/actions/workflows/ci.yml/badge.svg)](https://github.com/ezhigval/algo-sandbox/actions/workflows/ci.yml)
![License](https://img.shields.io/badge/license-MIT-blue)
![Tier](https://img.shields.io/badge/tier-junior-1d76db)

**[English](README.md)** · Русский

Рабочие реализации структур данных, о которых на собесах спрашивают снова и снова.

HTTP API на `:8084` и небольшой CLI для офлайн-демо.

## Зачем этот репозиторий

Нужно было одно место, где можно **показать** LRU, rate limiter'ы, consistent hashing и поиск по графу — а не только рассказывать на интервью.

## CLI

```bash
go run ./cmd/algo bench-lru
go run ./cmd/algo graph-demo
```

## HTTP (примеры)

```bash
make docker-up

# LRU в памяти
curl -X POST localhost:8084/api/v1/lru/put -d '{"key":"a","value":"1"}'
curl localhost:8084/api/v1/lru/get?key=a

# Top-K
curl -X POST localhost:8084/api/v1/topk -d '{"k":2,"items":[{"value":"a","score":1},{"value":"b","score":9}]}'

# BFS по графу
curl -X POST localhost:8084/api/v1/graph/path -d '{"from":"A","to":"D","graph":{"A":["B","C"],"B":["D"],"C":["D"],"D":[]}}'
```

Таблица сложностей: [COMPLEXITY.ru.md](COMPLEXITY.ru.md)

## Модули

LRU · token bucket · sliding window · consistent hash · top-K heap · BFS/DFS · worker pool

MIT
