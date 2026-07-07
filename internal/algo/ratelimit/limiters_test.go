package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Burst(t *testing.T) {
	b := NewTokenBucket(1, 2)
	assert.True(t, b.Allow())
	assert.True(t, b.Allow())
	assert.False(t, b.Allow())
}

func TestSlidingWindow(t *testing.T) {
	s := NewSlidingWindow(2, 100*time.Millisecond)
	assert.True(t, s.Allow())
	assert.True(t, s.Allow())
	assert.False(t, s.Allow())
	time.Sleep(110 * time.Millisecond)
	assert.True(t, s.Allow())
}

func BenchmarkTokenBucket(b *testing.B) {
	tb := NewTokenBucket(1000, 1000)
	for i := 0; i < b.N; i++ {
		tb.Allow()
	}
}
