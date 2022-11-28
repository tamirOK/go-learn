package cache

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

type DoublyLinkedList struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l *DoublyLinkedList) Len() int {
	return l.len
}

func (l *DoublyLinkedList) Front() *ListItem {
	return l.front
}

func (l *DoublyLinkedList) Back() *ListItem {
	return l.back
}

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

func (l *DoublyLinkedList) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	l.Remove(i)
	l.pushFront(i)
}
