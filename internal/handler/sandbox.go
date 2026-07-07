package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/ezhigval/go-toolkit/httputil"
	"github.com/ezhigval/algo-sandbox/internal/algo/consistent"
	"github.com/ezhigval/algo-sandbox/internal/algo/graph"
	"github.com/ezhigval/algo-sandbox/internal/algo/lru"
	"github.com/ezhigval/algo-sandbox/internal/algo/pool"
	"github.com/ezhigval/algo-sandbox/internal/algo/ratelimit"
	"github.com/ezhigval/algo-sandbox/internal/algo/topk"
	"github.com/redis/go-redis/v9"
)

type Sandbox struct {
	lruMem   *lru.Cache
	lruRedis *lru.RedisCache
	rdb      *redis.Client

	tbMu sync.Mutex
	tb   map[string]*ratelimit.TokenBucket

	swMu sync.Mutex
	sw   map[string]*ratelimit.SlidingWindow
}

func NewSandbox(rdb *redis.Client) *Sandbox {
	var rc *lru.RedisCache
	if rdb != nil {
		rc = lru.NewRedis(rdb, 100, time.Hour)
	}
	return &Sandbox{
		lruMem:   lru.New(100),
		lruRedis: rc,
		rdb:      rdb,
		tb:       make(map[string]*ratelimit.TokenBucket),
		sw:       make(map[string]*ratelimit.SlidingWindow),
	}
}

type lruPutReq struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Capacity int    `json:"capacity"`
	Backend  string `json:"backend"` // memory|redis
}

func (s *Sandbox) LRUPut(w http.ResponseWriter, r *http.Request) {
	var req lruPutReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, httputil.NewAppError(400, "BAD_REQUEST", "invalid json", err))
		return
	}
	if req.Backend == "redis" && s.lruRedis != nil {
		if err := s.lruRedis.Put(r.Context(), req.Key, req.Value); err != nil {
			httputil.WriteError(w, httputil.NewAppError(500, "REDIS_ERROR", "redis put failed", err))
			return
		}
	} else {
		if req.Capacity > 0 {
			s.lruMem = lru.New(req.Capacity)
		}
		s.lruMem.Put(req.Key, req.Value)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Sandbox) LRUGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	backend := r.URL.Query().Get("backend")

	if backend == "redis" && s.lruRedis != nil {
		val, err := s.lruRedis.Get(r.Context(), key)
		if err != nil {
			httputil.WriteError(w, httputil.NewAppError(500, "REDIS_ERROR", "redis get failed", err))
			return
		}
		if val == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		httputil.WriteJSON(w, http.StatusOK, map[string]string{"key": key, "value": val})
		return
	}

	val, ok := s.lruMem.Get(key)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"key": key, "value": val})
}

type rateReq struct {
	Key       string  `json:"key"`
	Rate      float64 `json:"rate"`
	Burst     int     `json:"burst"`
	Limit     int     `json:"limit"`
	WindowSec int     `json:"window_sec"`
}

func (s *Sandbox) TokenBucket(w http.ResponseWriter, r *http.Request) {
	var req rateReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Key == "" {
		req.Key = "default"
	}

	s.tbMu.Lock()
	b, ok := s.tb[req.Key]
	if !ok {
		b = ratelimit.NewTokenBucket(req.Rate, req.Burst)
		s.tb[req.Key] = b
	}
	s.tbMu.Unlock()

	httputil.WriteJSON(w, http.StatusOK, map[string]bool{"allowed": b.Allow()})
}

func (s *Sandbox) SlidingWindow(w http.ResponseWriter, r *http.Request) {
	var req rateReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	if req.Key == "" {
		req.Key = "default"
	}
	window := time.Duration(req.WindowSec) * time.Second
	if window <= 0 {
		window = time.Minute
	}

	s.swMu.Lock()
	sw, ok := s.sw[req.Key]
	if !ok {
		sw = ratelimit.NewSlidingWindow(req.Limit, window)
		s.sw[req.Key] = sw
	}
	s.swMu.Unlock()

	httputil.WriteJSON(w, http.StatusOK, map[string]bool{"allowed": sw.Allow()})
}

type hashReq struct {
	Key   string   `json:"key"`
	Nodes []string `json:"nodes"`
}

func (s *Sandbox) HashLookup(w http.ResponseWriter, r *http.Request) {
	var req hashReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, httputil.NewAppError(400, "BAD_REQUEST", "invalid json", err))
		return
	}
	ring := consistent.NewRing(req.Nodes, 50)
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"node": ring.Lookup(req.Key)})
}

type topKReq struct {
	K     int         `json:"k"`
	Items []topk.Item `json:"items"`
}

func (s *Sandbox) TopK(w http.ResponseWriter, r *http.Request) {
	var req topKReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, httputil.NewAppError(400, "BAD_REQUEST", "invalid json", err))
		return
	}
	httputil.WriteJSON(w, http.StatusOK, topk.TopK(req.Items, req.K))
}

type graphReq struct {
	Graph graph.Graph `json:"graph"`
	From  string      `json:"from"`
	To    string      `json:"to"`
	Algo  string      `json:"algo"`
}

func (s *Sandbox) GraphPath(w http.ResponseWriter, r *http.Request) {
	var req graphReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, httputil.NewAppError(400, "BAD_REQUEST", "invalid json", err))
		return
	}

	var path []string
	var ok bool
	switch req.Algo {
	case "dfs":
		path, ok = req.Graph.DFS(req.From, req.To)
	default:
		path, ok = req.Graph.BFS(req.From, req.To)
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"path": path, "algo": req.Algo})
}

type poolReq struct {
	Inputs  []int `json:"inputs"`
	Workers int   `json:"workers"`
}

func (s *Sandbox) PoolRun(w http.ResponseWriter, r *http.Request) {
	var req poolReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, httputil.NewAppError(400, "BAD_REQUEST", "invalid json", err))
		return
	}
	if req.Workers < 1 {
		req.Workers = 2
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"results": pool.RunInts(req.Inputs, req.Workers)})
}
