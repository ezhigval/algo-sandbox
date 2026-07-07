package ratelimit

import (
	"sync"
	"time"
)

// TokenBucket allows burst traffic with steady refill rate.
type TokenBucket struct {
	mu         sync.Mutex
	rate       float64 // tokens per second
	capacity   float64
	tokens     float64
	lastRefill time.Time
}

func NewTokenBucket(ratePerSec float64, burst int) *TokenBucket {
	if burst < 1 {
		burst = 1
	}
	return &TokenBucket{
		rate:       ratePerSec,
		capacity:   float64(burst),
		tokens:     float64(burst),
		lastRefill: time.Now(),
	}
}

func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.lastRefill = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// SlidingWindow counts events in a rolling time window (in-memory).
type SlidingWindow struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	events []time.Time
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{limit: limit, window: window}
}

func (s *SlidingWindow) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.window)

	// drop old
	i := 0
	for _, t := range s.events {
		if t.After(cutoff) {
			s.events[i] = t
			i++
		}
	}
	s.events = s.events[:i]

	if len(s.events) >= s.limit {
		return false
	}

	s.events = append(s.events, now)
	return true
}
