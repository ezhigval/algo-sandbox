package lru

import "container/list"

type Cache struct {
	cap   int
	items map[string]*list.Element
	order *list.List
}

type entry struct {
	key   string
	value string
}

func New(capacity int) *Cache {
	if capacity < 1 {
		capacity = 1
	}
	return &Cache{
		cap:   capacity,
		items: make(map[string]*list.Element),
		order: list.New(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	el, ok := c.items[key]
	if !ok {
		return "", false
	}
	c.order.MoveToFront(el)
	return el.Value.(*entry).value, true
}

func (c *Cache) Put(key, value string) {
	if el, ok := c.items[key]; ok {
		el.Value.(*entry).value = value
		c.order.MoveToFront(el)
		return
	}

	el := c.order.PushFront(&entry{key: key, value: value})
	c.items[key] = el

	if c.order.Len() > c.cap {
		back := c.order.Back()
		if back == nil {
			return
		}
		c.order.Remove(back)
		e := back.Value.(*entry)
		delete(c.items, e.key)
	}
}

func (c *Cache) Len() int { return c.order.Len() }
