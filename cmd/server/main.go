package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ezhigval/go-toolkit/httputil"
	"github.com/ezhigval/go-toolkit/logger"
	tkmw "github.com/ezhigval/go-toolkit/middleware"
	tkredis "github.com/ezhigval/go-toolkit/redis"
	"github.com/ezhigval/algo-sandbox/internal/config"
	"github.com/ezhigval/algo-sandbox/internal/handler"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(logger.Config{Level: cfg.LogLevel, Format: cfg.LogFormat})
	ctx := context.Background()

	var rdb *redis.Client
	if cfg.RedisAddr != "" {
		rdb = tkredis.NewClient(tkredis.Config{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
		defer func() { _ = tkredis.Close(rdb) }()
		if err := tkredis.Ping(ctx, rdb); err != nil {
			log.Warn("redis unavailable, lru redis backend disabled", "error", err)
			rdb = nil
		}
	}

	sb := handler.NewSandbox(rdb)

	r := chi.NewRouter()
	r.Use(tkmw.RequestID, tkmw.RealIP, tkmw.Recoverer(log), tkmw.AccessLog(log))
	r.Use(chimw.Timeout(30 * time.Second))

	r.Get("/health", httputil.HealthHandler(map[string]func() error{
		"redis": func() error {
			if rdb == nil {
				return nil
			}
			return tkredis.Ping(ctx, rdb)
		},
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/lru/put", sb.LRUPut)
		r.Get("/lru/get", sb.LRUGet)
		r.Post("/ratelimit/token-bucket", sb.TokenBucket)
		r.Post("/ratelimit/sliding-window", sb.SlidingWindow)
		r.Post("/hash/lookup", sb.HashLookup)
		r.Post("/topk", sb.TopK)
		r.Post("/graph/path", sb.GraphPath)
		r.Post("/pool/run", sb.PoolRun)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Info("algo sandbox listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
