package datastruct

import (
	"testing"
)

func TestList(t *testing.T) {
	list, err := ListCreate()
	if err != nil {
		t.Error("create list fail")
	}
	list.ListAddNodeTail(1)
	list.ListAddNodeTail(2)
	list.ListAddNodeTail(3)
	list.ListAddNodeTail(4)
	list.ListAddNodeTail("A")
	list.ListAddNodeTail("B")
	list.ListAddNodeTail("C")
	list.ListAddNodeTail("ABC123")

	length := list.ListLength()
	if length != 8 {
		t.Errorf("list length expect 8, actual %d", length)
	}

	iterator := list.ListGetIterator(AL_START_HEAD)
	next := ListNext(iterator)
	if next.value != 1 {
		t.Errorf("expect 1, actual %d", next.value)
	}

	ListNext(iterator)
	ListNext(iterator)
	ListNext(iterator)
	next = ListNext(iterator)
	if next.value != "A" {
		t.Errorf("expect A, actual %d", next.value)
	}

	ListNext(iterator)
	ListNext(iterator)
	next = ListNext(iterator)
	if next.value != "ABC123" {
		t.Errorf("expect ABC123, actual %d", next.value)
	}

	next = ListNext(iterator)
	if next != nil {
		t.Errorf("expect nil, actual %d", next)
	}

	index := list.ListIndex(3)
	if index.value != 4 {
		t.Errorf("expect 3, actual %d", index.value)
	}
}
