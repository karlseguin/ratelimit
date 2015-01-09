package ratelimit

type List struct {
	head *CachedTracker
	tail *CachedTracker
}

func NewList() *List {
	head := &CachedTracker{key: ""}
	tail := &CachedTracker{key: ""}
	head.next, tail.prev = tail, head
	return &List{
		head: head,
		tail: tail,
	}
}

func (l *List) PushToFront(item *CachedTracker) {
	l.Remove(item)
	head := l.head
	next := head.next
	next.prev = item
	item.next = next
	item.prev = head
	head.next = item
}

func (l *List) Remove(item *CachedTracker) {
	if item.prev == nil {
		return
	}
	item.prev.next, item.next.prev = item.next, item.prev
}
