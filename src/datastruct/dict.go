/**
字典hash实现
*/
package datastruct

import "time"

const (
	// 字典操作成功
	DICT_OK = iota
	// 字典操作失败
	DICT_ERR
)

// 哈希表节点的值
type dictValue struct {
	val *interface{}
	u64 uint64
	s64 int64
}

// 哈希表节点
type DictEntry struct {
	// 键
	key *interface{}
	// 值
	v dictValue
	// 指向下一个哈希表节点，hash冲突时生成链表使用
	next *DictEntry
}

// 字典类型特定函数
type dictType struct {
	// 计算哈希值的函数
	hashFunction func(key *interface{}) uint32
	// 复制键的函数
	keyDup func(privdata *interface{}, key *interface{}) *interface{}
	// 复制值的函数
	valDup func(privdata *interface{}, obj *interface{}) *interface{}
	// 对比键的函数
	keyCompare func(privdata *interface{}, key1 *interface{}, key2 *interface{}) int
	// 销毁键的函数
	keyDestructor func(privdata *interface{}, key *interface{})
	// 销毁值的函数
	valDestructor func(privdata *interface{}, obj *interface{})
}

// 哈希表
// 每个字典都使用两个哈希表，从而实现渐进式 rehash（从一个哈希表渐进转到另一个哈希表）
type dictht struct {
	// 哈希表数组
	table []*DictEntry
	// 哈希表大小
	size uint32
	// 哈希表大小掩码，用于计算索引值，总是等于 size -1
	// 一般使用hash值 & (size-1) 来定位索引
	sizemask uint32
	// 该哈希表已有节点的数量
	used uint32
}

// 字典
type dict struct {
	// 类型特定函数
	dtype *dictType
	// 私有数据
	privdata *interface{}
	// 哈希表
	ht [2]dictht
	// rehash进行到的索引，用于控制rehash进程
	// 当 rehash 不在进行时，值为-1
	rehshidx int
	// 目前正在运行的安全迭代器数量
	iterators int
}

// 字典迭代器
// 如果 safe 属性的值为 1 ，那么在迭代进行的过程中，
// 程序仍然可以执行 dictAdd 、 dictFind 和其他函数，对字典进行修改。
//
// 如果 safe 不为 1 ，那么程序只会调用 dictNext 对字典进行迭代，
// 而不对字典进行修改。
type dictIterator struct {
	// 被迭代的字典
	d *dict
	// 正在别迭代的哈希吗编号，可以是第0个或第1个
	table int
	// 迭代器当前所指向的哈希表索引位置
	index int
	// 标识这个迭代器是否安全 1:安全  0:不安全
	safe int
	// 当前迭代到的节点的指针
	entry *DictEntry
	// 当前迭代节点的下一个节点，
	// 因为在安全迭代器运作时，entry 所指向的节点可能会被修改，
	// 所以需要一个额外的指针来保存下一节点的位置，从而防止指针丢失
	nextEntry *DictEntry

	fingerprint int64
}

// 哈希表的数组的初始大小
const DICT_HT_INITIAL_SIZE uint32 = 4

// 释放给定字典节点的值
func (d *dict) dictFreeVal(entry *DictEntry) {
	valDestructor := d.dtype.valDestructor
	if valDestructor != nil {
		valDestructor(d.privdata, entry.v.val)
	}
}

// 设置给定字典节点的值
func (d *dict) dictSetVal(entry *DictEntry, val *interface{}) {
	if d.dtype.valDup != nil {
		entry.v.val = d.dtype.valDup(d.privdata, val)
	} else {
		entry.v.val = val
	}
}

// 将一个有符号整数设为节点的值
func dictSetSignedIntegerVal() {

}

// 将一个无符号整数设为节点的值
func dictSetUnsignedIntegerVal() {

}

// 释放给定字典节点的键
func dictFreeKey() {

}

// 设置给定字典节点的键
func dictSetKey(d *dict, entry *DictEntry, key *interface{}) {
	if d.dtype.keyDup != nil {
		entry.key = d.dtype.keyDup(d.privdata, key)
	} else {
		entry.key = key
	}
}

// 比对两个键
func dictCompareKeys(d *dict, key1 *interface{}, key2 *interface{}) bool {
	if d.dtype.keyCompare != nil {
		compare := d.dtype.keyCompare(d.privdata, key1, key2)
		return compare == 0
	}
	return key1 == key2
}

// 计算给定键的哈希值
func dictHashKey(d *dict, key *interface{}) uint32 {
	return d.dtype.hashFunction(key)
}

// 返回获取给定节点的键
func dictGetKey() {

}

// 返回获取给定节点的值
func dictGetVal() {

}

// 返回获取给定节点的有符号整数值
func dictGetSignedIntegerVal() {

}

// 返回给定节点的无符号整数值
func dictGetUnsignedIntegerVal() {

}

// 返回给定字典的大小
func dictSlots() {

}

// 返回字典的已有节点数量
func dictSize() {

}

// 查看字典是否正在 rehash
func dictIsRehashing(d *dict) bool {
	return d.rehshidx != -1
}

//==============================

// 指示字典是否启用rehash的标识
const dict_can_resize = 1

// 强制 rehash 的比率  即 已使用节点的数量 / 字典大小
const dict_force_resize_ratio = 5

// hash 函数
func DictIntHashFunction(key uint32) uint32 {
	key += ^(key << 15)
	key ^= (key >> 10)
	key += (key << 3)
	key ^= (key >> 6)
	key += ^(key << 11)
	key ^= (key >> 16)
	return key
}

var dict_hash_function_seed uint32 = 5381

func DictSetHashFunctionSeed(seed uint32) {
	dict_hash_function_seed = seed
}

func DictGetHashFunctionSeed() uint32 {
	return dict_hash_function_seed
}

func DictGenHashFunction(key interface{}, len int) uint32 {
	//因为Go无法操作指针移动，暂时无法实现C的代码
	return key.(uint32)
}

func DictGenCaseHashFunction(buf string) uint32 {
	var hash uint32 = 0
	for _, v := range buf {
		hash = hash*31 + uint32(v)
	}
	return hash
}

//================================ API implementation ====================

//基本不需要此函数
func dictReset(ht *dictht) {
	ht.table = nil
	ht.size = 0
	ht.sizemask = 0
	ht.used = 0
}

func DictCreate(dtype *dictType, privDataPtr *interface{}) *dict {
	d := &dict{}
	d.dictInit(dtype, privDataPtr)
	return d
}

// 初始化哈希表
func (d *dict) dictInit(dtype *dictType, privDataPtr *interface{}) int {
	d.ht[0] = dictht{}
	d.ht[1] = dictht{}

	d.dtype = dtype
	d.privdata = privDataPtr
	d.rehshidx = -1
	d.iterators = 0
	return DICT_OK
}

// 缩小给定字典
func (d *dict) dictResize() int {
	if dict_can_resize != 1 || dictIsRehashing(d) {
		return DICT_ERR
	}
	minimal := d.ht[0].used
	// 设置最小值
	if minimal < DICT_HT_INITIAL_SIZE {
		minimal = DICT_HT_INITIAL_SIZE
	}
	// 调整字典的大小
	return d.dictExpand(minimal)
}

// 创建一个新的哈希表，重新哈希到两个字典中未使用的字典，打开字典rehash标识
func (d *dict) dictExpand(size uint32) int {
	realSize := dictNextPower(size)

	// 不能在字典进行rehash 或size小于当前已使用节点时进行
	if dictIsRehashing(d) || d.ht[0].used > size {
		return DICT_ERR
	}

	// new hashtable
	n := dictht{}
	n.size = realSize
	n.sizemask = realSize - 1
	n.table = make([]*DictEntry, realSize)
	n.used = 0

	// 0号哈希表为空则说明没有填充过数据，这时进行初始化
	if d.ht[0].table == nil {
		d.ht[0] = n
		return DICT_OK
	}
	// 如果0号哈希表为空，则开始进行渐进式rehash
	d.ht[1] = n
	d.rehshidx = 0
	return DICT_OK
}

// 执行渐进式rehash
// 返回1表示未结束，仍有需要从0号哈希表移动到1号哈希表的数据
// n 为要进行rehash的数组数量(排除空元素)
func (d *dict) dictRehash(n int) int {
	// 只可以在开始了 rehash 开关后进行
	if !dictIsRehashing(d) {
		return 0
	}
	for n > 0 {
		n--
		// 0号哈希表为空，表示rehash结束,交换0号和1号哈希表
		if d.ht[0].used == 0 {
			d.ht[0] = d.ht[1]
			dictReset(&d.ht[1])
			d.rehshidx = -1
			return 0
		}
		// 跳过空数组
		for d.ht[0].table[d.rehshidx] == nil {
			d.rehshidx++
		}
		// 开始进行rehash
		de := d.ht[0].table[d.rehshidx]
		for de != nil {
			// 暂存下一个位置的地址
			nextde := de.next
			newIndex := dictHashKey(d, de.key) & d.ht[1].sizemask

			// 使用头插法插入到新的哈希表中的头部
			de.next = d.ht[1].table[newIndex]
			d.ht[1].table[newIndex] = de

			d.ht[0].used--
			d.ht[1].used++

			de = nextde
		}

		d.ht[0].table[d.rehshidx] = nil
		d.rehshidx++
	}
	return 1
}

// 返回以毫秒为单位的 UNIX 时间戳
func timeInMilliseconds() int64 {
	return time.Now().Unix() * 1000
}

// 在给定的毫秒内，已100步为不常，对字典进行rehash
func (d *dict) dictRehashMilliseconds(ms int) int {
	start := timeInMilliseconds()
	rehashes := 0
	for d.dictRehash(100) != 0 {
		rehashes += 100
		if timeInMilliseconds()-start > int64(ms) {
			break
		}
	}
	return rehashes
}

// 单步Rehash(没有安全迭代器的情况下)
func (d *dict) dictRehashStep() {
	if d.iterators == 0 {
		d.dictRehash(1)
	}
}

// 将key,value 添加到字典中
func (d *dict) dictAdd(key *interface{}, val *interface{}) int {
	entry := d.dictAddRaw(key)
	if entry == nil {
		return DICT_ERR
	}
	d.dictSetVal(entry, val)
	return DICT_OK
}

// 将key插入到字典中(不包括值)
func (d *dict) dictAddRaw(key *interface{}) *DictEntry {
	// 如果渐进hash在进行中，那么在新增时执行一次单步rehash
	if dictIsRehashing(d) {
		d.dictRehashStep()
	}
	index := dictKeyIndex(d, key)
	// -1 表示键已经存在
	if index == -1 {
		return nil
	}
	var ht *dictht
	if dictIsRehashing(d) {
		ht = &d.ht[1]
	} else {
		ht = &d.ht[0]
	}
	entry := &DictEntry{}
	// 头插法插入链表节点
	entry.next = ht.table[index]
	ht.table[index] = entry
	ht.used++

	dictSetKey(d, entry, key)
	return entry
}

// 添加、替换 key-value到dict中
func (d *dict) dictReplace(key *interface{}, val *interface{}) int {
	// 能新增则新增后返回
	if d.dictAdd(key, val) == DICT_OK {
		return 1
	}
	entry := dictFind(d, key)
	d.dictSetVal(entry, val)
	return 0
}

// 添加key到dict,如果已经存在则直接返回
func (d *dict) dictReplaceRaw(key *interface{}) *DictEntry {
	dictEntry := dictFind(d, key)
	if dictEntry != nil {
		return d.dictAddRaw(key)
	}
	return dictEntry
}

// 返回字典表中包含key的节点，查询不到返回nil
func dictFind(d *dict, key *interface{}) *DictEntry {
	// 0号哈希表为空则表示整个dict为空
	if d.ht[0].sizemask == 0 {
		return nil
	}
	// 如果在rehash过程中，则单步执行一步
	if dictIsRehashing(d) {
		d.dictRehashStep()
	}

	hash := dictHashKey(d, key)
	for table := 0; table <= 1; table++ {
		idx := hash & d.ht[table].sizemask
		he := d.ht[table].table[idx]
		for he != nil {
			if dictCompareKeys(d, key, he.key) {
				return he
			}
			he = he.next
		}
		// 没在进行rehash则说明1号哈希表为空，不需搜索
		if !dictIsRehashing(d) {
			return nil
		}
	}
	return nil
}

// 计算key的索引，如果已经存在，返回-1
// 计算时需要考虑是否在渐进rehash进程中，来决定是插入到哪个哈希表
// 进行中插入到1号哈希表，否则插入到0号哈希表
func dictKeyIndex(d *dict, key *interface{}) int32 {
	if dictExpandIfNeeded(d) == DICT_ERR {
		return -1
	}
	hash := dictHashKey(d, key)
	var idx int32
	for table := 0; table <= 1; table++ {
		idx = int32(hash & d.ht[table].sizemask)
		// 判断相同的key是否已存在
		// 定位到索引后，从链表(如果存在)往下找
		he := d.ht[table].table[idx]
		for he != nil {
			if dictCompareKeys(d, key, he.key) {
				return -1
			}
			he = he.next
		}
		// 如果0号哈希表中没有这个key，并且没有在进行rehash
		// 此时以0号哈希表中的索引为准，结束循环
		if !dictIsRehashing(d) {
			break
		}
	}
	return idx
}

// 初始化字典或者满足条件时进行扩展
func dictExpandIfNeeded(d *dict) int {
	// 如果已经在rehash过程中，直接返回
	if dictIsRehashing(d) {
		return DICT_OK
	}
	if d.ht[0].size == 0 {
		return d.dictExpand(DICT_HT_INITIAL_SIZE)
	}

	if d.ht[0].used >= d.ht[0].size &&
		(dict_can_resize == 1 || d.ht[0].used/d.ht[0].size > dict_force_resize_ratio) {
		return d.dictExpand(d.ht[0].used * 2)
	}
	return DICT_OK

}

// 计算第一个大于等于size的2的N次方
func dictNextPower(size uint32) uint32 {
	i := DICT_HT_INITIAL_SIZE
	if size >= LONG_MAX {
		return LONG_MAX + 1
	}
	for i < size {
		i *= 2
	}
	return i
}
