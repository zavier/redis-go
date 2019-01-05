package datastruct

import (
	"unsafe"
)

const ZSKPLIST_MAXLEVEL = 32
const ZSKIPLIST_P = 0.25

// redis对象类型
const (
	REDIS_STRING = iota
	REDIS_LIST
	REDIS_SET
	REDIS_ZSET
	REDIS_HASH
)

// 对象编码
const (
	REDIS_ENCODING_RAW = iota // sds
	REDIS_ENCODING_INT
	REDIS_ENCODING_HT // hash table
	REDIS_ENCODING_ZIPMAP
	REDIS_ENCODING_LINKEDLIST
	REDIS_ENCODING_ZIPLIST
	REDIS_ENCODING_INTSET
	REDIS_ENCODING_SKIPLIST
	REDIS_ENCODING_EMBSTR // embeded string encoding
)

func sdsEncodedObject(objptr *redisObject) bool {
	return objptr.encoding == REDIS_ENCODING_RAW || objptr.encoding == REDIS_ENCODING_EMBSTR
}

//======================================================

// Redis对象
type redisObject struct {
	// 类型
	rtype byte
	// 编码
	encoding byte
	// 对象最后一次被访问的时间
	lru uint32
	// 引用计数
	refcount int
	// 指向实际值的指针
	ptr unsafe.Pointer
}

// 跳跃表节点
type zskiplistNode struct {
	// 成员对象-redisObject
	obj *redisObject
	// 分值
	score float64
	// 后退指针
	backward *zskiplistNode
	// 层
	level []zskiplistLevel
}

// 层
type zskiplistLevel struct {
	// 前进指针
	forward *zskiplistNode
	// 距离同层下一个节点跨度(两个链表节点之前的距离)
	span int
}

// 跳跃表
type zskiplist struct {
	// 表头、表尾节点指针
	header, tail *zskiplistNode
	// 表中节点数量
	length int
	// 表中层数最大的节点层数
	level int
}
