package datastruct

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
	ptr *interface{}
}

// 跳跃表节点
type zskiplistNode struct {
	// 成员对象
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
	// 跨度
	span uint32
}
