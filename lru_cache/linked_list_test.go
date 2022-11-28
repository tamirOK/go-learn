package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoublyLinkedListLen(t *testing.T) {
	list := DoublyLinkedList{}

	assert.Equal(t, list.Len(), 0)

	for i := 0; i < 10; i++ {
		list.PushFront(i)
	}

	assert.Equal(t, list.Len(), 10)
}

func TestDoublyLinkedListFront(t *testing.T) {
	list := DoublyLinkedList{}

	assert.Nil(t, list.Front())

	list.PushFront(0)

	assert.Equal(t, list.Front().Value, 0)
	assert.Nil(t, list.Front().Prev)

	list.PushBack(1)

	assert.Equal(t, list.Front().Value, 0)
	assert.Nil(t, list.Front().Prev)
}

func TestDoublyLinkedListBack(t *testing.T) {
	list := DoublyLinkedList{}

	assert.Nil(t, list.Back())

	list.PushBack(0)

	assert.Equal(t, list.Back().Value, 0)
	assert.Nil(t, list.Back().Next)

	list.PushFront(1)

	assert.Equal(t, list.Back().Value, 0)
	assert.Nil(t, list.Back().Next)
}

func TestDoublyLinkedListPushFront(t *testing.T) {
	list := DoublyLinkedList{}

	for i := 0; i < 5; i++ {
		list.PushFront(i)

		if i == 0 {
			// ensure front and back are same, when list contains 1 item
			assert.Equal(t, list.Front(), list.Back())
		} else {
			assert.NotEqual(t, list.Front(), list.Back())
		}

		assert.Equal(t, list.Front().Value, i)
		assert.Equal(t, list.Len(), i+1)
	}

	// ensure Prev updated correctly on each PushFront operation
	for value, item := 0, list.Back(); item != nil; item, value = item.Prev, value+1 {
		assert.Equal(t, item.Value, value)
	}
}

func TestDoublyLinkedListPushBack(t *testing.T) {
	list := DoublyLinkedList{}

	for i := 0; i < 5; i++ {
		list.PushBack(i)

		if i == 0 {
			// ensure front and back are same, when list contains 1 item
			assert.Equal(t, list.Front(), list.Back())
		} else {
			assert.NotEqual(t, list.Front(), list.Back())
		}

		assert.Equal(t, list.Back().Value, i)
		assert.Equal(t, list.Len(), i+1)
	}

	// ensure Next updated correctly on each PushBack operation
	for value, item := 0, list.Front(); item != nil; item, value = item.Next, value+1 {
		assert.Equal(t, item.Value, value)
	}
}

func TestDoublyLinkedListRemove(t *testing.T) {
	list := DoublyLinkedList{}

	one := list.PushBack(1)
	two := list.PushFront(2)
	three := list.PushBack(3)
	// list: 2 <-> 1 <-> 3

	assert.Equal(t, list.Len(), 3)
	assert.Equal(t, two.Next, one)
	assert.Equal(t, three.Prev, one)

	list.Remove(one)
	// list: 2 <-> 3

	assert.Equal(t, list.Len(), 2)
	assert.Equal(t, two.Next, three)
	assert.Equal(t, three.Prev, two)
}

func TestDoublyLinkedListRemoveFront(t *testing.T) {
	list := DoublyLinkedList{}

	one := list.PushBack(1)
	two := list.PushFront(2)
	list.PushBack(3)
	// list: 2 <-> 1 <-> 3

	assert.Equal(t, list.Front(), two)

	list.Remove(two)
	// list: 1 <-> 3

	assert.Equal(t, list.Front(), one)
	assert.Nil(t, list.Front().Prev)
}

func TestDoublyLinkedListRemoveBack(t *testing.T) {
	list := DoublyLinkedList{}

	one := list.PushBack(1)
	list.PushFront(2)
	three := list.PushBack(3)
	// list: 2 <-> 1 <-> 3

	assert.Equal(t, list.Back(), three)

	list.Remove(three)
	// list: 2 <-> 1

	assert.Equal(t, list.Back(), one)
	assert.Nil(t, list.Back().Next)
}

func TestDoublyLinkedListRemoveSingleItem(t *testing.T) {
	list := DoublyLinkedList{}

	one := list.PushBack(1)
	list.Remove(one)

	assert.Equal(t, list.Len(), 0)
	assert.Nil(t, list.Back())
	assert.Nil(t, list.Front())
}

func TestDoublyLinkedListMoveToFront(t *testing.T) {
	list := DoublyLinkedList{}

	one := list.PushBack(1)
	list.PushFront(2)
	list.PushBack(3)
	// list: 2 <-> 1 <-> 3

	list.MoveToFront(one)
	// list: 1 <-> 2 <-> 3

	assert.Equal(t, list.Len(), 3)
	assert.Equal(t, list.Front(), one)
}

func TestDoublyLinkedListMoveToFrontFrontItem(t *testing.T) {
	list := DoublyLinkedList{}

	list.PushBack(1)
	two := list.PushFront(2)
	list.PushBack(3)
	// list: 2 <-> 1 <-> 3

	list.MoveToFront(two)
	// list: 2 <-> 1 <-> 3

	assert.Equal(t, list.Len(), 3)
	assert.Equal(t, list.Front(), two)
}

func TestDoublyLinkedListMoveToFrontBackItem(t *testing.T) {
	list := DoublyLinkedList{}

	one := list.PushBack(1)
	list.PushFront(2)
	three := list.PushBack(3)
	// list: 2 <-> 1 <-> 3

	list.MoveToFront(three)
	// list: 3 <-> 2 <-> 1

	assert.Equal(t, list.Len(), 3)
	assert.Equal(t, list.Front(), three)
	assert.Equal(t, list.Back(), one)
}
