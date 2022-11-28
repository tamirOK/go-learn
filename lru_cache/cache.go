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

type LRUCache struct {
	capacity int
	queue    DoublyLinkedList
	storage  map[Key]*ListItem
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		queue:    DoublyLinkedList{},
		storage:  make(map[Key]*ListItem, capacity),
	}
}

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

func (c *LRUCache) Get(key Key) (interface{}, bool) {
	listItem, ok := c.storage[key]

	if !ok {
		return nil, false
	}

	cachedItem := listItem.Value.(*cacheItem)
	c.queue.MoveToFront(listItem)

	return cachedItem.value, true
}

func (c *LRUCache) Clear() {
	c.queue = DoublyLinkedList{}
	c.storage = make(map[Key]*ListItem, c.capacity)
}
