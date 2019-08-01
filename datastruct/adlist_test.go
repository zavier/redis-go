package datastruct

import (
	"testing"
)

func TestListCreate(t *testing.T) {
	list, err := ListCreate()
	if err != nil || list.ListLength() != 0 || list.head != nil {
		t.Error("create list fail")
	}
}

func TestListEmpty(t *testing.T) {
	list, err := ListCreate()
	if err != nil || list.ListLength() != 0 || list.head != nil {
		t.Error("create list fail")
		return
	}
	list.ListAddNodeTail(1)
	ListEmpty(list)
	if list.ListLength() != 0 || list.head != nil {
		t.Error("create list fail")
	}
}

func TestList_ListAddNodeHead(t *testing.T) {
	list, _ := ListCreate()
	list.ListAddNodeHead("123")
	list.ListAddNodeHead("456")
	list.ListAddNodeHead("abc")
	if list.ListLength() != 3 {
		t.Error("ListAddNodeHead length error")
	}
	res := checkList(list, []string{"abc", "456", "123"})
	if !res {
		t.Error("ListAddNodeHead error")
	}
}

func TestList_ListAddNodeTail(t *testing.T) {
	list, _ := ListCreate()
	list.ListAddNodeTail("123")
	list.ListAddNodeTail("456")
	list.ListAddNodeTail("abc")
	if list.ListLength() != 3 {
		t.Error("ListAddNodeHead length error")
	}
	res := checkList(list, []string{"123", "456", "abc"})
	if !res {
		t.Error("ListAddNodeHead error")
	}
}

func TestList_ListInsertNode(t *testing.T) {
	list, _ := ListCreate()
	list.ListAddNodeTail("123")
	list.ListAddNodeTail("456")
	list.ListAddNodeTail("abc")

	node := list.head
	list.ListInsertNode(node, "wqeurp", 0)
	res := checkList(list, []string{"wqeurp", "123", "456", "abc"})
	if !res || list.ListLength() != 4 {
		t.Error("ListInsertNode error")
	}

	node = list.head
	list.ListInsertNode(node, "name", 1)
	res = checkList(list, []string{"wqeurp", "name", "123", "456", "abc"})
	if !res || list.ListLength() != 5 {
		t.Error("ListInsertNode error")
	}
}

func TestList_ListDelNode(t *testing.T) {
	list, _ := ListCreate()
	list.ListAddNodeTail("123")
	list.ListAddNodeTail("456")
	list.ListAddNodeTail("abc")

	list.ListDelNode(list.head.next)
	res := checkList(list, []string{"123", "abc"})
	if !res || list.ListLength() != 2 {
		t.Error("ListDelNode error")
	}
}

// 校验list中的元素与切片中的内容是否一致
func checkList(list *List, expec []string) bool {
	node := list.head
	for _, item := range expec {
		if node.value != item {
			return false
		}
		node = node.next
	}
	if node != nil {
		return false
	}
	return true
}

func TestList_ListGetIterator(t *testing.T) {
	list, _ := ListCreate()
	iterator := list.ListGetIterator(AL_START_HEAD)
	if iterator == nil {
		t.Error("ListGetIterator error")
	}
}

func TestList_ListRewind(t *testing.T) {
	list, _ := ListCreate()
	list.ListAddNodeTail("123")
	list.ListAddNodeTail("456")
	iterator := list.ListGetIterator(AL_START_HEAD)
	list.ListRewind(iterator)
	if iterator.next.value != list.ListGetIterator(AL_START_HEAD).next.value {
		t.Error("ListRewind error")
	}
}

func TestList_ListRewindTail(t *testing.T) {
	list, _ := ListCreate()
	list.ListAddNodeTail("123")
	list.ListAddNodeTail("456")
	iterator := list.ListGetIterator(AL_START_TAIL)
	list.ListRewindTail(iterator)
	if iterator.next.value != list.ListGetIterator(AL_START_TAIL).next.value {
		t.Error("ListRewind error")
	}
}

func TestListNode_ListNext(t *testing.T) {
	list, _ := ListCreate()
	strList := []string{"123", "456", "789", "abc"}

	for _, item := range strList {
		list.ListAddNodeTail(item)
	}

	iterator := list.ListGetIterator(AL_START_HEAD)
	i := 0
	for node := ListNext(iterator); node != nil; node = ListNext(iterator) {
		if node.value != strList[i] {
			t.Error("ListNext error")
		}
		i++
	}

	iterator = list.ListGetIterator(AL_START_TAIL)
	i = 3
	for node := ListNext(iterator); node != nil; node = ListNext(iterator) {
		if node.value != strList[i] {
			t.Error("ListNext error")
		}
		i--
	}
}

func TestList_ListDup(t *testing.T) {
	list, _ := ListCreate()
	strList := []string{"123", "456", "789", "abc"}
	for _, item := range strList {
		list.ListAddNodeTail(item)
	}

	dup := list.ListDup()
	res := checkList(dup, strList)
	if !res {
		t.Error("ListDup error")
		return
	}

	dup.head.value = "321"
	if list.head.value == "321" {
		t.Error("ListDup error")
		return
	}
}

func TestList_ListSearchKey(t *testing.T) {
	list, _ := ListCreate()
	strList := []string{"123", "456", "789", "abc"}
	for _, item := range strList {
		list.ListAddNodeTail(item)
	}

	node := list.ListSearchKey("456")
	if node.value != "456" {
		t.Error("ListSearchKey error")
	}
	if node.prev.value != "123" {
		t.Error("ListSearchKey error")
	}
	if node.next.value != "789" {
		t.Error("ListSearchKey error")
	}
}

func TestList_ListIndex(t *testing.T) {
	list, _ := ListCreate()
	strList := []string{"123", "456", "789", "abc"}
	for _, item := range strList {
		list.ListAddNodeTail(item)
	}

	node := list.ListIndex(0)
	if node.value != "123" {
		t.Error("ListIndex error")
	}

	node = list.ListIndex(-1)
	if node.value != "abc" {
		t.Error("ListIndex error")
	}

	node = list.ListIndex(-3)
	if node.value != "456" {
		t.Error("ListIndex error")
	}
}

func TestList_ListRotate(t *testing.T) {
	list, _ := ListCreate()
	strList := []string{"123", "456", "789", "abc"}
	for _, item := range strList {
		list.ListAddNodeTail(item)
	}

	list.ListRotate()

	checkList(list, []string{"abc", "123", "456", "789"})
}

func TestList_ListJoin(t *testing.T) {
	list1, _ := ListCreate()
	list1.ListAddNodeTail("1")
	list1.ListAddNodeTail("1")
	list1.ListAddNodeTail("1")

	list2, _ := ListCreate()
	list2.ListAddNodeTail("2")
	list2.ListAddNodeTail("2")

	list1.ListJoin(list2)

	res := checkList(list1, []string{"1", "1", "1", "2", "2"})
	if !res || list1.ListLength() != 5 {
		t.Error("ListJoin error")
	}

	if list2.head != nil || list2.len != 0 {
		t.Error("ListJoin error")
	}
}
