package datastruct

import "math/rand"

// 创建一个层数为level的跳跃表节点
// 成员对象为 obj, 分值为 score
func zslCreateNode(level int, score float64, obj *redisObject) *zskiplistNode {
	znode := &zskiplistNode{}
	znode.score = score
	znode.obj = obj
	return znode
}

// 创建并返回一个新的跳跃表
func zslCreate() *zskiplist {
	zsl := &zskiplist{}
	zsl.level = 1
	zsl.length = 0

	zsl.header = zslCreateNode(ZSKPLIST_MAXLEVEL, 0, nil)
	for j := 0; j < ZSKPLIST_MAXLEVEL; j++ {
		zsl.header.level[j].forward = nil
		zsl.header.level[j].span = 0
	}
	zsl.header.backward = nil
	zsl.tail = nil
	return zsl
}

// 释放给定的跳跃表节点
func zslFreeNode(node *zskiplistNode) {
	decrRefCount(node.obj)
	node = nil
}

// 释放给定跳跃表，以及表中的所有节点
func zslFree(zsl *zskiplist) {
	node := zsl.header.level[0].forward
	var next *zskiplistNode
	zsl.header = nil
	for node != nil {
		next = node.level[0].forward
		zslFreeNode(node)
		node = next
	}
	zsl = nil
}

// 返回一个随机值 [1, ZSKIPLIST_MAXLEVEL]
// 用作新跳跃表节点的层数
func zslRandomLevel() int {
	level := 1
	i := float64(rand.Int() & 0xffff)
	f := ZSKIPLIST_P * 0xffff
	for i < f {
		i = float64(rand.Int() & 0xffff)
		level++
	}
	if level < ZSKPLIST_MAXLEVEL {
		return level
	}
	return ZSKPLIST_MAXLEVEL
}

// 创建一个新节点，并插入跳跃表中
func zslInsert(zsl *zskiplist, score float64, robj *redisObject) *zskiplistNode {
	var rank [ZSKPLIST_MAXLEVEL]int
	var update []*zskiplistNode
	x := zsl.header
	for i := zsl.level - 1; i >= 0; i-- {
		if i == zsl.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		forwareNode := x.level[i].forward
		for forwareNode != nil {
			if forwareNode.score < score ||
				(forwareNode.score == score &&
					compareStringObjects(x.level[i].forward.obj, robj) < 0) {
				// 记录跨越过了多少节点
				rank[i] += x.level[i].span
				// 移动至下一个指针
				x = x.level[i].forward
			} else {
				break
			}
		}
		update[i] = x
	}

	level := zslRandomLevel()
	if level > zsl.level {
		for i := zsl.level; i < level; i++ {
			rank[i] = 0
			update[i] = zsl.header
			update[i].level[i].span = zsl.length
		}
		zsl.level = level
	}

	x = zslCreateNode(level, score, robj)
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := level; i < zsl.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == zsl.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		zsl.tail = x
	}

	zsl.length++
	return x
}

// 内部删除函数
func zslDeleteNode() {
	//todo
}
