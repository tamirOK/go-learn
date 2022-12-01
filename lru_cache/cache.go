// Package cache implements LRU cache based on doubly linked list.
package cache

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) (bool, error)
	Get(key Key) (interface{}, bool)
	Clear()
}

// LRUCache is an implementation of a cache with the least recently used policy.
// It implements Cache interface.
type LRUCache struct {
	capacity int
	queue    DoublyLinkedList
	storage  map[Key]*ListItem
}

// NewLRUCache returns pointer to newly created LRUCache type with given capacity.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		queue:    DoublyLinkedList{},
		storage:  make(map[Key]*ListItem, capacity),
	}
}

// Set stores key with given value.
// Sets return boolean indicating existence of a given key in the cache.
func (c *LRUCache) Set(key Key, value interface{}) bool {
	listItem, ok := c.storage[key]

	// if element with given key already in cache
	if ok {
		cachedItem := listItem.Value.(*cacheItem)
		cachedItem.value = value
		c.queue.MoveToFront(listItem)

		return true
	}

	// if cache is full
	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		cachedItem := lastItem.Value.(*cacheItem)
		c.queue.Remove(lastItem)
		delete(c.storage, cachedItem.key)
	}

	c.storage[key] = c.queue.PushFront(
		&cacheItem{
			key:   key,
			value: value,
		},
	)

	return false
}

// Get returns value associated with a given key in the cache and
// boolean indicating existence of key in the cache.
func (c *LRUCache) Get(key Key) (interface{}, bool) {
	listItem, ok := c.storage[key]

	if !ok {
		return nil, false
	}

	cachedItem := listItem.Value.(*cacheItem)
	c.queue.MoveToFront(listItem)

	return cachedItem.value, true
}

// Clear removes all data in the cache.
func (c *LRUCache) Clear() {
	c.queue = DoublyLinkedList{}
	c.storage = make(map[Key]*ListItem, c.capacity)
}
