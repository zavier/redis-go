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

/*
双端链表迭代器
*/
type ListIter struct {
	// 当前迭代到的节点
	next *ListNode
	// 迭代的方向
	direction int
}

/*
双端链表结构
*/
type List struct {
	// 表头节点
	head *ListNode
	// 表尾节点
	tail *ListNode
	// 节点值复制函数
	dup func(value NodeValue) NodeValue
	// 节点值释放函数
	free func()
	// 节点值对比函数
	match func(value1 NodeValue, value2 NodeValue) bool
	// 链表所包含的节点数量
	len int
}

// 返回链表所包含的节点数量
func (self *List) ListLength() int {
	return self.len
}

func (self *List) ListFirst() *ListNode {
	return self.head
}

func (self *List) ListLast() *ListNode {
	return self.tail
}

type NodeValue interface{}

/*
双端链表节点
*/
type ListNode struct {
	prev  *ListNode
	next  *ListNode
	value NodeValue
}

func (self *ListNode) ListPrevNode() *ListNode {
	return self.prev
}

func (self *ListNode) ListNextNode() *ListNode {
	return self.next
}

func (self *ListNode) ListNodeValue() NodeValue {
	return self.value
}

// 设置链表的值复制函数为f
func (self *List) ListSetDupMethod(f func()) {
	self.dup = f
}

// 设置链表的值释放函数为f
func (self *List) ListSetFreeMethod(f func()) {
	self.free = f
}

// 设置链表的对比函数为f
func (self *List) ListSetMatchMethod(f func(value1 NodeValue, value2 NodeValue) bool) {
	self.match = f
}

// 返回链表的值复制函数
func (self *List) ListGetDupMethod() func(value NodeValue) NodeValue {
	return self.dup
}

// 返回链表的值释放函数
func (self *List) ListGetFree() func() {
	return self.free
}

// 返回链表的值对比函数
func (self *List) ListGetMatchMethod() func(value1 NodeValue, value2 NodeValue) bool {
	return self.match
}

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
func (self *List) ListAddNodeHead(value NodeValue) {
	node := &ListNode{}
	node.value = value
	if self.len == 0 {
		self.head, self.tail = node, node
		node.prev, node.next = nil, nil
	} else {
		node.prev = nil
		node.next = self.head
		self.head.prev = node
		self.head = node
	}
	self.len++
}

// 添加节点到表尾，尾插法
func (self *List) ListAddNodeTail(value NodeValue) {
	node := &ListNode{}
	node.value = value
	if self.len == 0 {
		self.head, self.tail = node, node
		node.prev, node.next = nil, nil
	} else {
		node.prev = self.tail
		node.next = nil
		self.tail.next = node
		self.tail = node
	}
	self.len++
}

// 创建一个新节点，将其添加到 oldNode 节点之前或之后
// 如果 after 为 0，添加到 oldNode 节点之前
// 如果 after 为 1，添加到 oldNode 节点之后
func (self *List) ListInsertNode(oldNode *ListNode, value NodeValue, after int) {
	node := &ListNode{}
	node.value = value
	// 添加到给定节点之后
	if after == 1 {
		node.prev = oldNode
		node.next = oldNode.next
		// 如果oldNode原为表尾节点
		if self.tail == oldNode {
			self.tail = node
		}
	} else {
		// 添加节点到指定节点之前
		node.next = oldNode
		node.prev = oldNode.prev
		// 如果给定节点是表头节点
		if self.head == oldNode {
			self.head = node
		}
	}

	// 更新新节点的前后节点的对应指针指向当前节点
	if node.prev != nil {
		node.prev.next = node
	}
	if node.next != nil {
		node.next.prev = node
	}

	self.len++
}

// 删除指定节点
func (self *List) ListDelNode(node *ListNode) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		self.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		self.tail = node.prev
	}

	self.len--
}

// 为链表创建一个迭代器
// direction  AL_START_HEAD ：从表头向表尾迭代
// direction  AL_START_TAIL ：从表尾向表头迭代
func (self *List) ListGetIterator(direction int) *ListIter {
	iter := &ListIter{}
	if direction == AL_START_HEAD {
		iter.next = self.head
	} else {
		iter.next = self.tail
	}
	iter.direction = direction
	return iter
}

// 释放迭代器
func ListReleaseIterator(iter *ListIter) {
	iter = nil
}

// 将迭代器的方向设置为 AL_START_HEAD
// 并将迭代器的指针重新指向表头节点
func (self *List) ListRewind(li *ListIter) {
	li.next = self.head
	li.direction = AL_START_HEAD
}

// 将迭代器的方向设置为 AL_START_TAIL
// 并将迭代器的指针重新指向表尾节点
func (self *List) ListRewindTail(li *ListIter) {
	li.next = self.tail
	li.direction = AL_START_TAIL
}

// 返回迭代器当前所指向的节点
func ListNext(iter *ListIter) *ListNode {
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

func (self *List) ListDup() *List {
	list, err := ListCreate()
	if err != nil {
		return nil
	}
	list.dup = self.dup
	list.free = self.free
	list.match = self.match

	iter := self.ListGetIterator(AL_START_HEAD)
	node := ListNext(iter)
	for node != nil {
		var value NodeValue
		// 如果有复制函数，则使用复制函数复制值
		if list.dup != nil {
			value = list.dup(node.value)
			if value == nil {
				ListRelease(list)
				ListReleaseIterator(iter)
				return nil
			}
		} else {
			value = node.value
		}
		// 将节点添加到链表
		list.ListAddNodeTail(value)

		node = ListNext(iter)
	}
	// 释放迭代器
	ListReleaseIterator(iter)
	return list
}

// 查询链表中的key值节点
// 对比操作由链表的 match 函数负责进行
// 如果没有设置 match 函数，则直接比较值
// 匹配成功返回第一个匹配的节点，否则返回nil
func (self *List) ListSearchKey(key NodeValue) *ListNode {
	iter := self.ListGetIterator(AL_START_HEAD)
	node := ListNext(iter)
	for node != nil {
		if self.match != nil {
			if self.match(node.value, key) {
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
func (self *List) ListIndex(index int) *ListNode {
	var n *ListNode
	// 如果索引为负数，从表尾开始查找
	if index < 0 {
		index = (-index) - 1
		n = self.tail
		for index > 0 && n != nil {
			n = n.prev
			index--
		}
	} else {
		n = self.head
		for index > 0 && n != nil {
			n = n.next
			index--
		}
	}
	return n
}

// 取出链表的表尾节点，并将它移动到表头，成为新的表头节点
func (self *List) ListRotate() {
	tail := self.tail
	if self.ListLength() <= 1 {
		return
	}

	// 重置表尾节点
	self.tail = tail.prev
	self.tail.next = nil
	// 重置表头节点
	self.head.prev = tail
	tail.prev = nil
	tail.next = self.head
	self.head = tail
}
