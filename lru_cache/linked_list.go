package cache

// ListItem represents single item of a list.
type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

// DoublyLinkedList is the implementation of doulby linked list data structure.
// It implements List interface.
type DoublyLinkedList struct {
	len   int
	front *ListItem
	back  *ListItem
}

// Len returns a length of the doubly linked list.
func (l *DoublyLinkedList) Len() int {
	return l.len
}

// Front returns first element of the doubly linked list.
func (l *DoublyLinkedList) Front() *ListItem {
	return l.front
}

// Back returns last element of the doubly linked list.
func (l *DoublyLinkedList) Back() *ListItem {
	return l.back
}

// PushFront puts given value to the beginning of the doubly linked list.
// it can be very fast and sometimes it is very slow.
func (l *DoublyLinkedList) PushFront(v interface{}) *ListItem {
	newFront := &ListItem{
		Value: v,
	}
	l.pushFront(newFront)

	return newFront
}

func (l *DoublyLinkedList) pushFront(newFront *ListItem) {
	newFront.Next = l.front
	newFront.Prev = nil

	if l.front != nil {
		l.front.Prev = newFront
	}
	l.front = newFront

	// create tail if list is empty
	if l.back == nil {
		l.back = newFront
	}

	l.len++
}

// PushFront puts given value to the end of the doubly linked list.
func (l *DoublyLinkedList) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{
		Value: v,
		Prev:  l.back,
	}
	if l.back != nil {
		l.back.Next = newBack
	}
	l.back = newBack

	// create head if list is empty
	if l.front == nil {
		l.front = newBack
	}

	l.len++

	return newBack
}

// Remove removes given item from the doubly linked list.
func (l *DoublyLinkedList) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	l.len--
}

// MoveToFront moves given item to the beginnin of the doubly linked list.
func (l *DoublyLinkedList) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	l.Remove(i)
	l.pushFront(i)
}
