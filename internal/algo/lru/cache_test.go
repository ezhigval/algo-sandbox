package lru

import (
	"fmt"
	"testing"
)

func TestCache_EvictsLRU(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Put("b", "2")
	c.Get("a")
	c.Put("c", "3")

	if _, ok := c.Get("b"); ok {
		t.Fatal("b should be evicted")
	}
	if v, ok := c.Get("a"); !ok || v != "1" {
		t.Fatal("a should remain")
	}
}

func BenchmarkCachePutGet(b *testing.B) {
	c := New(128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("k%d", i%256)
		c.Put(key, key)
		c.Get(key)
	}
}
