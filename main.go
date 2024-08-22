package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	capacity int
	items    map[string]*list.Element
	queue    *list.List
	mutex    sync.RWMutex
}

type entry struct {
	key   string
	value interface{}
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		queue:    list.New(),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.queue.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.queue.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	if c.queue.Len() >= c.capacity {
		oldest := c.queue.Back()
		if oldest != nil {
			c.queue.Remove(oldest)
			delete(c.items, oldest.Value.(*entry).key)
		}
	}

	elem := c.queue.PushFront(&entry{key, value})
	c.items[key] = elem
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.queue.Remove(elem)
		delete(c.items, key)
	}
}

func (c *Cache) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.queue.Len()
}
