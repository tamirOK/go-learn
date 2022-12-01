package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertQueueFrontValue(t *testing.T, c *LRUCache, expectedValue string) {
	t.Helper()
	cachedItem := c.queue.Front().Value.(*cacheItem)

	assert.Equal(t, cachedItem.value, expectedValue)
}

func assertKeyDoesNotExist(t *testing.T, c *LRUCache, key Key) {
	t.Helper()
	value, ok := c.Get(key)

	assert.Nil(t, value)
	assert.False(t, ok)
}

func TestLRUCacheSetResultWithNewElement(t *testing.T) {
	c := NewLRUCache(5)

	assert.False(t, c.Set("testKey", "testValue"))
}

func TestLRUCacheSetResultWithExistingElement(t *testing.T) {
	c := NewLRUCache(5)
	c.Set("testKey", "testValue")

	assert.True(t, c.Set("testKey", "newTestValue"))

	result, ok := c.Get("testKey")

	assert.Equal(t, result, "newTestValue")
	assert.True(t, ok)
}

func TestLRUCacheEnsureElementMovesFrontWithSet(t *testing.T) {
	c := NewLRUCache(2)
	c.Set("testKey1", "testValue1")

	assertQueueFrontValue(t, c, "testValue1")

	c.Set("testKey2", "testValue2")

	assertQueueFrontValue(t, c, "testValue2")

	c.Set("testKey3", "testValue3")

	assertQueueFrontValue(t, c, "testValue3")
}

func TestLRUCacheSetWithFullCapacity(t *testing.T) {
	c := NewLRUCache(2)
	c.Set("testKey1", "testValue1")
	c.Set("testKey2", "testValue2")
	c.Set("testKey3", "testValue3")

	assertKeyDoesNotExist(t, c, "testKey1")

	c.Set("testKey2", "newTestValue2")
	c.Set("testKey4", "testValue4")

	assertKeyDoesNotExist(t, c, "testKey3")
}

func TestLRUCacheGet(t *testing.T) {
	c := NewLRUCache(5)
	c.Set("testKey", "testValue")
	value, ok := c.Get("testKey")

	assert.Equal(t, value, "testValue")
	assert.True(t, ok)
}

func TestLRUCacheEnsureElementMovesFrontWithGet(t *testing.T) {
	c := NewLRUCache(5)
	c.Set("testKey1", "testValue1")
	c.Set("testKey2", "testValue2")

	c.Get("testKey1")

	assertQueueFrontValue(t, c, "testValue1")

	c.Get("testKey2")

	assertQueueFrontValue(t, c, "testValue2")

	c.Get("NonExistentKey")

	assertQueueFrontValue(t, c, "testValue2")
}

func TestLRUCacheGetNonexistentKey(t *testing.T) {
	c := NewLRUCache(5)

	assertKeyDoesNotExist(t, c, "nonexistentKey")
}

func TestLRUCacheClear(t *testing.T) {
	c := NewLRUCache(5)
	c.Set("testKey", "testValue")
	c.Clear()

	assertKeyDoesNotExist(t, c, "testKey")
}
