package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer func() { l.len++ }()
	item := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.Len() == 0 {
		l.addFirstItem(item)
		return item
	}

	l.front.Prev = item
	l.front = item

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	defer func() { l.len++ }()
	item := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.Len() == 0 {
		l.addFirstItem(item)
		return item
	}

	l.back.Next = item
	l.back = item

	return item
}

func (l *list) Remove(i *ListItem) {
	defer func() { l.len-- }()

	if i.Next == nil && i.Prev == nil {
		l.front = nil
		l.back = nil
		return
	}

	if i.Next == nil {
		i.Prev.Next = nil
		l.back = i.Prev
		return
	}

	if i.Prev == nil {
		i.Next.Prev = nil
		l.front = i.Next
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)

	l.front.Prev = i

	i.Next = l.front
	i.Prev = nil

	l.front = i
	l.len++ //.Remove compensation
}

func (l *list) addFirstItem(i *ListItem) {
	l.front = i
	l.back = i
}

func NewList() List {
	return new(list)
}
