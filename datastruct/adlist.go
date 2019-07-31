/*
Redis(3.0)中双端链表实现的简单仿写
*/
package datastruct

const (
	// 从表头向表尾进行迭代
	AL_START_HEAD = iota
	// 从表尾到表头进行迭代
	AL_START_TAIL
)

// 双端链表节点
type listNode struct {
	// 前置节点
	prev *listNode
	// 后置节点
	next *listNode
	// 节点的值
	value interface{}
}

// 双端链表迭代器
type listIter struct {
	// 当前迭代到的节点
	next *listNode
	// 迭代的方向
	direction int
}

// 双端链表结构
type List struct {
	// 表头节点
	head *listNode
	// 表尾节点
	tail *listNode
	// 节点值复制函数
	dup func(value interface{}) interface{}
	// 节点值释放函数
	free func()
	// 节点值对比函数
	match func(value1 interface{}, value2 interface{}) bool
	// 链表所包含的节点数量
	len int
}

// 返回链表所包含的节点数量
func (list *List) ListLength() int {
	return list.len
}

// 返回给定链表的表头节点
func (list *List) ListFirst() *listNode {
	return list.head
}

// 返回给定链表的表尾节点
func (list *List) ListLast() *listNode {
	return list.tail
}

// 返回给定节点的前置节点
func (list *listNode) ListPrevNode() *listNode {
	return list.prev
}

// 返回给定节点的后置位置
func (list *listNode) ListNextNode() *listNode {
	return list.next
}

// 返回给定节点的值
func (list *listNode) ListNodeValue() interface{} {
	return list.value
}

// 设置链表的值复制函数为f
func (list *List) ListSetDupMethod(f func(value interface{}) interface{}) {
	list.dup = f
}

// 设置链表的值释放函数为f
func (list *List) ListSetFreeMethod(f func()) {
	list.free = f
}

// 设置链表的对比函数为f
func (list *List) ListSetMatchMethod(f func(value1 interface{}, value2 interface{}) bool) {
	list.match = f
}

// 返回链表的值复制函数
func (list *List) ListGetDupMethod() func(value interface{}) interface{} {
	return list.dup
}

// 返回链表的值释放函数
func (list *List) ListGetFree() func() {
	return list.free
}

// 返回链表的值对比函数
func (list *List) ListGetMatchMethod() func(value1 interface{}, value2 interface{}) bool {
	return list.match
}

//===========================================================

// 创建一个新的链表
func ListCreate() (list *List, err error) {
	list = &List{}
	list.head, list.tail = nil, nil
	list.dup, list.free = nil, nil
	list.len = 0
	list.match = nil
	return list, nil
}

// 释放整个链表，以及链表中的所有节点
func ListRelease(list *List) {
	list = nil
}

// 添加节点到表头，头插法
func (list *List) ListAddNodeHead(value interface{}) {
	node := &listNode{}
	node.value = value
	if list.len == 0 {
		list.head, list.tail = node, node
		node.prev, node.next = nil, nil
	} else {
		node.prev = nil
		node.next = list.head
		list.head.prev = node
		list.head = node
	}
	list.len++
}

// 添加节点到表尾，尾插法
func (list *List) ListAddNodeTail(value interface{}) {
	node := &listNode{}
	node.value = value
	if list.len == 0 {
		list.head, list.tail = node, node
		node.prev, node.next = nil, nil
	} else {
		node.prev = list.tail
		node.next = nil
		list.tail.next = node
		list.tail = node
	}
	list.len++
}

// 创建一个新节点，将其添加到 oldNode 节点之前或之后
// 如果 after 为 0，添加到 oldNode 节点之前
// 如果 after 为 1，添加到 oldNode 节点之后
func (list *List) ListInsertNode(oldNode *listNode, value interface{}, after int) {
	node := &listNode{}
	node.value = value
	// 添加到给定节点之后
	if after == 1 {
		node.prev = oldNode
		node.next = oldNode.next
		// 如果oldNode原为表尾节点
		if list.tail == oldNode {
			list.tail = node
		}
	} else {
		// 添加节点到指定节点之前
		node.next = oldNode
		node.prev = oldNode.prev
		// 如果给定节点是表头节点
		if list.head == oldNode {
			list.head = node
		}
	}

	// 更新新节点的前后节点的对应指针指向当前节点
	if node.prev != nil {
		node.prev.next = node
	}
	if node.next != nil {
		node.next.prev = node
	}

	list.len++
}

// 删除指定节点
func (list *List) ListDelNode(node *listNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		list.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		list.tail = node.prev
	}

	list.len--
}

// 为链表创建一个迭代器
// direction  AL_START_HEAD ：从表头向表尾迭代
// direction  AL_START_TAIL ：从表尾向表头迭代
func (list *List) ListGetIterator(direction int) *listIter {
	iter := &listIter{}
	if direction == AL_START_HEAD {
		iter.next = list.head
	} else {
		iter.next = list.tail
	}
	iter.direction = direction
	return iter
}

// 释放迭代器
func ListReleaseIterator(iter *listIter) {
	iter = nil
}

// 将迭代器的方向设置为 AL_START_HEAD
// 并将迭代器的指针重新指向表头节点
func (list *List) ListRewind(li *listIter) {
	li.next = list.head
	li.direction = AL_START_HEAD
}

// 将迭代器的方向设置为 AL_START_TAIL
// 并将迭代器的指针重新指向表尾节点
func (list *List) ListRewindTail(li *listIter) {
	li.next = list.tail
	li.direction = AL_START_TAIL
}

// 返回迭代器当前所指向的节点
func ListNext(iter *listIter) *listNode {
	current := iter.next
	if current != nil {
		if iter.direction == AL_START_HEAD {
			iter.next = current.next
		} else {
			iter.next = current.prev
		}
	}
	return current
}

// 复制整个链表
func (list *List) ListDup() *List {
	newList, err := ListCreate()
	if err != nil {
		return nil
	}
	newList.dup = list.dup
	newList.free = list.free
	newList.match = list.match

	iter := list.ListGetIterator(AL_START_HEAD)
	node := ListNext(iter)
	for node != nil {
		var value interface{}
		// 如果有复制函数，则使用复制函数复制值
		if newList.dup != nil {
			value = newList.dup(node.value)
			if value == nil {
				ListRelease(newList)
				ListReleaseIterator(iter)
				return nil
			}
		} else {
			value = node.value
		}
		// 将节点添加到链表
		newList.ListAddNodeTail(value)

		node = ListNext(iter)
	}
	// 释放迭代器
	ListReleaseIterator(iter)
	return newList
}

// 查询链表中的key值节点
// 对比操作由链表的 match 函数负责进行
// 如果没有设置 match 函数，则直接比较值
// 匹配成功返回第一个匹配的节点，否则返回nil
func (list *List) ListSearchKey(key interface{}) *listNode {
	iter := list.ListGetIterator(AL_START_HEAD)
	node := ListNext(iter)
	for node != nil {
		if list.match != nil {
			if list.match(node.value, key) {
				ListReleaseIterator(iter)
				return node
			}
		} else {
			if key == node.value {
				ListReleaseIterator(iter)
				return node
			}
		}
		node = ListNext(iter)
	}
	ListReleaseIterator(iter)
	return nil
}

// 返回链表在给定索引上的值
// 索引可以为负数，超出索引范围返回nil
func (list *List) ListIndex(index int) *listNode {
	var n *listNode
	// 如果索引为负数，从表尾开始查找
	if index < 0 {
		index = (-index) - 1
		n = list.tail
		for index > 0 && n != nil {
			n = n.prev
			index--
		}
	} else {
		n = list.head
		for index > 0 && n != nil {
			n = n.next
			index--
		}
	}
	return n
}

// 取出链表的表尾节点，并将它移动到表头，成为新的表头节点
func (list *List) ListRotate() {
	tail := list.tail
	if list.ListLength() <= 1 {
		return
	}

	// 重置表尾节点
	list.tail = tail.prev
	list.tail.next = nil
	// 重置表头节点
	list.head.prev = tail
	tail.prev = nil
	tail.next = list.head
	list.head = tail
}
