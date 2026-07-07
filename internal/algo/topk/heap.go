package topk

import "container/heap"

type Item struct {
	Value string
	Score int
}

type minHeap []Item

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].Score < h[j].Score }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x any) { *h = append(*h, x.(Item)) }
func (h *minHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// TopK returns k highest-scored items. O(n log k).
func TopK(items []Item, k int) []Item {
	if k <= 0 || len(items) == 0 {
		return nil
	}
	if k >= len(items) {
		out := append([]Item(nil), items...)
		sortDesc(out)
		return out
	}

	h := &minHeap{}
	heap.Init(h)

	for _, it := range items {
		if h.Len() < k {
			heap.Push(h, it)
			continue
		}
		if it.Score > (*h)[0].Score {
			heap.Pop(h)
			heap.Push(h, it)
		}
	}

	out := make([]Item, h.Len())
	for i := h.Len() - 1; i >= 0; i-- {
		out[i] = heap.Pop(h).(Item)
	}
	return out
}

func sortDesc(items []Item) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].Score > items[i].Score {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}
