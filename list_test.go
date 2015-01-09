package ratelimit

import (
	. "github.com/karlseguin/expect"
	"testing"
)

type ListTests struct{}

func Test_List(t *testing.T) {
	Expectify(new(ListTests), t)
}

func (_ ListTests) PushesNewItemToFront() {
	list := NewList()
	cr1, cr2, cr3 := &CachedTracker{key: "2"}, &CachedTracker{key: "1"}, &CachedTracker{key: "3"}
	list.PushToFront(cr2)
	list.PushToFront(cr3)
	list.PushToFront(cr1)
	assertList(list, cr1, cr3, cr2)
}

func (_ ListTests) PushesExistingItems() {
	list := NewList()
	cr1, cr2, cr3 := &CachedTracker{key: "2"}, &CachedTracker{key: "1"}, &CachedTracker{key: "3"}
	list.PushToFront(cr2)
	list.PushToFront(cr2)
	assertList(list, cr2)
	list.PushToFront(cr3)
	assertList(list, cr3, cr2)
	list.PushToFront(cr3)
	assertList(list, cr3, cr2)
	list.PushToFront(cr1)
	assertList(list, cr1, cr3, cr2)
	list.PushToFront(cr2)
	assertList(list, cr2, cr1, cr3)
}

func (_ ListTests) RemovesItemFromTheList() {
	list := NewList()
	cr1, cr2, cr3, cr4, cr5 := &CachedTracker{key: "2"}, &CachedTracker{key: "1"}, &CachedTracker{key: "3"}, &CachedTracker{key: "4"}, &CachedTracker{key: "5"}
	list.PushToFront(cr2)
	list.PushToFront(cr3)
	list.PushToFront(cr5)
	list.PushToFront(cr1)
	list.PushToFront(cr4)
	list.Remove(cr5)
	assertList(list, cr4, cr1, cr3, cr2)
	list.Remove(cr4)
	assertList(list, cr1, cr3, cr2)
	list.Remove(cr2)
	assertList(list, cr1, cr3)
	list.Remove(cr1)
	assertList(list, cr3)
	list.Remove(cr3)
	assertList(list)
}

func (_ ListTests) NoopOnRemovingNoExistingItem() {
	list := NewList()
	list.Remove(&CachedTracker{key: "x"})
	assertList(list)
}

func assertList(list *List, items ...*CachedTracker) {
	Expect(list.head.key).To.Equal("")
	Expect(list.tail.key).To.Equal("")
	node := list.head
	for _, item := range items {
		node = node.next
		Expect(item).To.Equal(node)
	}
	Expect(node.next.key).To.Equal("")
}
