package consistent

import (
	"hash/crc32"
	"sort"
	"sync"
)

type Ring struct {
	mu          sync.RWMutex
	virtualNodes int
	nodes       []string
	ring        []uint32
	hashToNode  map[uint32]string
}

func NewRing(nodes []string, virtualNodes int) *Ring {
	if virtualNodes < 1 {
		virtualNodes = 100
	}
	r := &Ring{
		virtualNodes: virtualNodes,
		nodes:        append([]string(nil), nodes...),
		hashToNode:   make(map[uint32]string),
	}
	r.rebuild()
	return r
}

func (r *Ring) Lookup(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.ring) == 0 {
		return ""
	}

	h := hash(key)
	idx := sort.Search(len(r.ring), func(i int) bool { return r.ring[i] >= h })
	if idx == len(r.ring) {
		idx = 0
	}
	return r.hashToNode[r.ring[idx]]
}

func (r *Ring) AddNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes = append(r.nodes, node)
	r.rebuild()
}

func (r *Ring) rebuild() {
	r.ring = r.ring[:0]
	r.hashToNode = make(map[uint32]string)

	for _, node := range r.nodes {
		for i := range r.virtualNodes {
			h := hash(node + "#" + itoa(i))
			r.ring = append(r.ring, h)
			r.hashToNode[h] = node
		}
	}
	sort.Slice(r.ring, func(i, j int) bool { return r.ring[i] < r.ring[j] })
}

func hash(s string) uint32 { return crc32.ChecksumIEEE([]byte(s)) }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
