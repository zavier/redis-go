/**
字典hash实现
*/
package datastruct

import (
	"errors"
	"math/rand"
	"time"
)

const (
	// 字典操作成功
	DICT_OK = iota
	// 字典操作失败
	DICT_ERR
)

// 哈希表节点的值
type dictValue struct {
	val interface{}
	u64 uint64
	s64 int64
}

// 哈希表节点
type DictEntry struct {
	// 键
	key interface{}
	// 值
	v dictValue
	// 指向下一个哈希表节点，hash冲突时生成链表使用
	next *DictEntry
}

// 字典类型特定函数
type dictType struct {
	// 计算哈希值的函数
	hashFunction func(key interface{}) int
	// 复制键的函数
	keyDup func(privdata interface{}, key interface{}) interface{}
	// 复制值的函数
	valDup func(privdata interface{}, obj interface{}) interface{}
	// 对比键的函数
	keyCompare func(privdata interface{}, key1 interface{}, key2 interface{}) int
	// 销毁键的函数
	keyDestructor func(privdata interface{}, key interface{})
	// 销毁值的函数
	valDestructor func(privdata interface{}, obj interface{})
}

// 哈希表
// 每个字典都使用两个哈希表，从而实现渐进式 rehash（从一个哈希表渐进转到另一个哈希表）
type dictht struct {
	// 哈希表数组
	table []*DictEntry
	// 哈希表大小
	size int
	// 哈希表大小掩码，用于计算索引值，总是等于 size -1
	// 一般使用hash值 & (size-1) 来定位索引
	sizemask int
	// 该哈希表已有节点的数量
	used int
}

// 字典
type dict struct {
	// 类型特定函数
	dtype dictType
	// 私有数据
	privdata interface{}
	// 哈希表(一共两个，rehash使用)
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
const DICT_HT_INITIAL_SIZE int = 4

// 释放给定字典节点的值
func (d *dict) dictFreeVal(entry DictEntry) {
	valDestructor := d.dtype.valDestructor
	if valDestructor != nil {
		valDestructor(d.privdata, entry.v.val)
	}
}

// 设置给定字典节点的值
func (d *dict) dictSetVal(entry *DictEntry, val interface{}) {
	if d.dtype.valDup != nil {
		entry.v.val = d.dtype.valDup(d.privdata, val)
	} else {
		entry.v.val = val
	}
}

// 将一个有符号整数设为节点的值
func dictSetSignedIntegerVal(entry *DictEntry, val interface{}) {
	entry.v.s64 = val.(int64)
}

// 将一个无符号整数设为节点的值
func dictSetUnsignedIntegerVal(entry *DictEntry, val interface{}) {
	entry.v.u64 = val.(uint64)
}

// 释放给定字典节点的键
func dictFreeKey(d *dict, entry *DictEntry) {
	if d.dtype.keyDestructor != nil {
		d.dtype.keyDestructor(d.privdata, entry.key)
	}
}

// 设置给定字典节点的键
func dictSetKey(d *dict, entry *DictEntry, key interface{}) {
	if d.dtype.keyDup != nil {
		entry.key = d.dtype.keyDup(d.privdata, key)
	} else {
		entry.key = key
	}
}

// 比对两个键
func dictCompareKeys(d *dict, key1 interface{}, key2 interface{}) bool {
	if d.dtype.keyCompare != nil {
		compare := d.dtype.keyCompare(d.privdata, key1, key2)
		return compare == 0
	}
	return key1 == key2
}

// 计算给定键的哈希值
func dictHashKey(d *dict, key interface{}) int {
	return d.dtype.hashFunction(key)
}

// 返回获取给定节点的键
func dictGetKey() {

}

// 返回获取给定节点的值
func dictGetVal(he *DictEntry) interface{} {
	return he.v.val
}

// 返回获取给定节点的有符号整数值
func dictGetSignedIntegerVal(he *DictEntry) int64 {
	return he.v.s64
}

// 返回给定节点的无符号整数值
func dictGetUnsignedIntegerVal(he *DictEntry) uint64 {
	return he.v.u64
}

// 返回给定字典的大小
func dictSlots(d *dict) int {
	return d.ht[0].size + d.ht[1].size
}

// 返回字典的已有节点数量
func dictSize(d *dict) int {
	return d.ht[0].size + d.ht[1].size
}

// 查看字典是否正在 rehash
func dictIsRehashing(d *dict) bool {
	return d.rehshidx != -1
}

//==============================

// 指示字典是否启用rehash的标识
var dict_can_resize = 1

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
	// todo
	return uint32(5)
}

func DictGenCaseHashFunction(buf string) uint32 {
	var hash uint32 = 0
	for _, v := range buf {
		hash = hash*31 + uint32(v)
	}
	return hash
}

//================================ API implementation ====================

// 重置哈希表
func dictReset(ht *dictht) {
	ht.table = nil
	ht.size = 0
	ht.sizemask = 0
	ht.used = 0
}

// 创建一个字典
func DictCreate(dtype dictType, privDataPtr interface{}) *dict {
	d := &dict{}
	d.dictInit(dtype, privDataPtr)
	return d
}

// 初始化哈希表
func (d *dict) dictInit(dtype dictType, privDataPtr interface{}) int {
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

// 创建一个新的哈希表，或者重新哈希到两个字典中未使用的字典，并打开字典rehash标识
func (d *dict) dictExpand(size int) int {
	realSize := dictNextPower(size)

	// 不能在字典进行rehash 或size小于当前已使用节点时进行
	if dictIsRehashing(d) || d.ht[0].used > size {
		return DICT_ERR
	}

	// new hashtable
	n := dictht{}
	n.size = realSize
	n.sizemask = realSize - 1
	n.table = make([]*DictEntry, 0, realSize)
	n.used = 0

	// 0号哈希表为空则说明没有填充过数据，这时进行初始化
	if d.ht[0].table == nil {
		d.ht[0] = n
		return DICT_OK
	}
	// 如果0号哈希表非空，则开始进行rehash
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
func (d *dict) dictAdd(key interface{}, val interface{}) int {
	entry := d.dictAddRaw(key)
	if entry == nil {
		return DICT_ERR
	}
	d.dictSetVal(entry, val)
	return DICT_OK
}

// 将key插入到字典中(不包括值)
func (d *dict) dictAddRaw(key interface{}) *DictEntry {
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
func (d *dict) dictReplace(key interface{}, val interface{}) int {
	// 能新增则新增后返回
	if d.dictAdd(key, val) == DICT_OK {
		return 1
	}
	entry := dictFind(d, key)
	d.dictSetVal(entry, val)
	return 0
}

// 添加key到dict,如果已经存在则直接返回
func (d *dict) dictReplaceRaw(key interface{}) *DictEntry {
	dictEntry := dictFind(d, key)
	if dictEntry != nil {
		return d.dictAddRaw(key)
	}
	return dictEntry
}

// 查找并删除包含给定键的节点
func dictGenericDelete(d *dict, key interface{}) int {
	if d.ht[0].size == 0 {
		return DICT_ERR
	}

	if dictIsRehashing(d) {
		d.dictRehashStep()
	}

	hash := dictHashKey(d, key)
	for table := 0; table <= 1; table++ {
		idx := hash & d.ht[table].sizemask

		// 拿到对应索引所在数组中的头结点
		he := d.ht[table].table[idx]
		var prevHe *DictEntry
		for he != nil {
			if dictCompareKeys(d, key, he.key) {
				// 头结点就是要找的key对应节点
				if prevHe == nil {
					d.ht[table].table[idx] = he.next
				} else {
					prevHe.next = he.next
				}

				d.ht[table].used--
				return DICT_OK
			}

			prevHe = he
			he = he.next
		}

		// 0号哈希表未找到，且没在进行哈希，也就不用查找1号哈希表了
		if !dictIsRehashing(d) {
			break
		}
	}

	return DICT_ERR
}

// 从字典中删除包含给定键的节点, 本来区分是是否调用释放函数，此处不提供不释放的实现
var dictDelete = dictGenericDelete
var dictDeleteNoFree = dictGenericDelete

// 删除哈希表上的所有节点，重置属性
func dictClear(d *dict, ht *dictht) int {
	dictReset(ht)
	return DICT_OK
}

// 删除并释放整个字典
func dictRelease(d *dict) {
	d.ht[0] = dictht{}
	d.ht[1] = dictht{}
}

// 返回字典表中包含key的节点，查询不到返回nil
func dictFind(d *dict, key interface{}) *DictEntry {
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

// 返回包含指定key的value值
func dictFetchValue(d *dict, key interface{}) interface{} {
	he := dictFind(d, key)
	if he == nil {
		return nil
	}
	return dictGetVal(he)
}

func dictFingerprint(d *dict) int64 {
	//todo
	return 1
}

// 创建并返回字典的不安全迭代器
func dictGetIterator(d *dict) *dictIterator {
	iterator := &dictIterator{}
	iterator.d = d
	iterator.table = 0
	iterator.index = -1
	iterator.safe = 0
	iterator.entry = nil
	iterator.nextEntry = nil
	return iterator
}

// 创建并返回给定节点的安全迭代器
func dictGetSafeIterator(d *dict) *dictIterator {
	iterator := dictGetIterator(d)
	iterator.safe = 1
	return iterator
}

// 返回迭代器指向的当前节点，结束返回nil
func dictNext(iter *dictIterator) *DictEntry {
	for {
		if iter.entry == nil {
			// 获取哈希表
			ht := iter.d.ht[iter.table]
			// 如果是第一次迭代
			if iter.index == -1 && iter.table == 0 {
				// 安全则更新安全迭代器计数器
				if iter.safe == 1 {
					iter.d.iterators++
				} else {
					iter.fingerprint = dictFingerprint(iter.d)
				}
			}
			iter.index++
			//
			if iter.index >= int(ht.size) {
				// 如果在进行rehash,则迭代1号哈希表
				if dictIsRehashing(iter.d) && iter.table == 0 {
					iter.table++
					iter.index = 0
					ht = iter.d.ht[1]
				} else {
					break
				}
			}

			iter.entry = ht.table[iter.index]
		} else {
			iter.entry = iter.nextEntry
		}

		if iter.entry != nil {
			iter.nextEntry = iter.entry.next
			return iter.entry
		}
	}
	return nil
}

// 释放给定字典迭代器
func dictReleaseIterator(iter *dictIterator) {
	// 已经开始进行迭代
	if !(iter.index == -1 && iter.table == 0) {
		if iter.safe == 1 {
			iter.d.iterators--
		} else {
			if iter.fingerprint != dictFingerprint(iter.d) {
				panic(errors.New("fingerprint not equals"))
			}
		}
	}
	iter = nil
}

// 随机返回一个节点
func dictGetRandomKey(d *dict) *DictEntry {
	if dictSize(d) == 0 {
		return nil
	}
	if dictIsRehashing(d) {
		d.dictRehashStep()
	}

	var he *DictEntry
	if dictIsRehashing(d) {
		h := rand.Intn(int(d.ht[0].size + d.ht[1].size))
		for he == nil {
			if h > d.ht[0].size {
				he = d.ht[1].table[h-d.ht[0].size]
			} else {
				he = d.ht[0].table[h]
			}
		}
	} else {
		for he == nil {
			h := uint(rand.Intn(int(d.ht[0].sizemask + 1)))
			he = d.ht[0].table[h]
		}
	}

	// 获取当前数组节点拥有的链表长度
	listlen := 0
	orighe := he
	for he != nil {
		he = he.next
		listlen++
	}
	// 获取链表上的随机一个数据
	listlen = rand.Intn(listlen)
	he = orighe
	for listlen > 0 {
		listlen--
		he = he.next
	}
	return he
}

// 随机返回count个节点，可能有重复数据
func dictGetRandomKeys(d *dict, count int) ([]*DictEntry, int) {
	if dictSize(d) < count {
		count = int(dictSize(d))
	}
	stored := 0
	dest := make([]*DictEntry, 0, count)
	for stored < count {
		for j := 0; j < 2; j++ {
			// 随机获取一个索引
			i := rand.Intn(d.ht[j].size)
			size := d.ht[j].size

			for size > 0 {
				size--
				he := d.ht[j].table[i]
				// 将这个索引对应的链表都获取到
				for he != nil {
					dest = append(dest, he)
					he = he.next
					stored++
					if stored == count {
						return dest, stored
					}
				}

				// 如果加上索引i中的数据不够，则在随机获取下一个索引
				i = (i + 1) & d.ht[j].sizemask
			}
			// 到这里表示0号哈希表数据不够，这时候dict必须处于rehash中，这时去1号哈希表找
			if !dictIsRehashing(d) {
				panic(errors.New("must be rehashing"))
			}
		}
	}
	return dest, stored
}

// 翻转位 from: http://graphics.stanford.edu/~seander/bithacks.html#ReverseParallel
func rev(v uint32) uint32 {
	// todo
	return v
}

type dictScanFunction struct {
	privdata interface{}
	de       *DictEntry
}

// todo 改写的有些问题
func dictScan(d *dict, v int, fn *dictScanFunction, privdata interface{}) int {
	if dictSize(d) == 0 {
		return 0
	}
	var t0, t1 *dictht
	var de *DictEntry
	var m0, m1 int
	if !dictIsRehashing(d) {
		t0 = &d.ht[0]
		m0 = t0.sizemask

		de = t0.table[v&m0]
		for de != nil {
			fn = &dictScanFunction{privdata, de}
			de = de.next
		}
	} else {
		t0, t1 = &d.ht[0], &d.ht[1]
		// 确保t0的size小于t1
		if t0.size > t1.size {
			t0, t1 = t1, t0
		}
		m0, m1 = t0.sizemask, t1.sizemask

		de = t0.table[v&m0]
		for de != nil {
			fn = &dictScanFunction{privdata, de}
			de = de.next
		}

		res := 1
		for res > 0 {
			de = t1.table[v&m1]
			for de != nil {
				fn = &dictScanFunction{privdata, de}
				de = de.next
			}
			v = (((v | m0) + 1) & ^m0) | (v & m0)
			res = v & (m0 ^ m1)
		}
	}

	v |= ^m0
	v = int(rev(uint32(v)))
	v++
	v = int(rev(uint32(v)))
	return v
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
func dictNextPower(size int) int {
	i := DICT_HT_INITIAL_SIZE
	if uint(size) >= LONG_MAX {
		return int(LONG_MAX)
	}
	for i < size {
		i *= 2
	}
	return i
}

// 计算key的索引，如果已经存在，返回-1
// 计算时需要考虑是否在渐进rehash进程中，来决定是插入到哪个哈希表
// 进行中插入到1号哈希表，否则插入到0号哈希表
func dictKeyIndex(d *dict, key interface{}) int32 {
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

// 清空所有哈希表节点
func dictEmpty(d *dict) {
	d.ht[0] = dictht{}
	d.ht[1] = dictht{}
	d.rehshidx = -1
	d.iterators = 0
}

// 开始自动rehash
func dictEnableResize() {
	dict_can_resize = 1
}

// 关闭自动rehash
func dictDisableResize() {
	dict_can_resize = 0
}
