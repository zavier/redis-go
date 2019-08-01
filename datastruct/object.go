package datastruct

import "errors"

const REDIS_COMPARE_BINARY = 1 << 0
const REDIS_COMPARE_COLL = 1 << 1

// 释放字符串对象
func freeStringObject(robj *redisObject) {
	if robj.encoding == REDIS_ENCODING_RAW {
		//ptr := robj.ptr
		//sdsFree((*ptr).(sds))
		robj.ptr = nil
	}
}

// 释放列表对象
func freeListObject(robj *redisObject) {
	switch robj.encoding {
	case REDIS_ENCODING_LINKEDLIST:
		robj.ptr = nil
	case REDIS_ENCODING_ZIPLIST:
		robj.ptr = nil
	default:
		panic(errors.New("Unknown list encoding type"))
	}
}

// 释放集合对象
func freeSetObject(robj *redisObject) {
	switch robj.encoding {
	case REDIS_ENCODING_HT:
		dictRelease((*dict)(robj.ptr))
	case REDIS_ENCODING_INTSET:
		robj.ptr = nil
	default:
		panic(errors.New("Unknown set encoding type"))
	}
}

// 释放有序集合对象
func freeZsetObject(robj *redisObject) {
	//todo
}

// 释放哈希对象
func freeHashObject(robj *redisObject) {
	switch robj.encoding {
	case REDIS_ENCODING_HT:
		dictRelease((*dict)(robj.ptr))
	case REDIS_ENCODING_ZIPLIST:
		robj.ptr = nil
	default:
		panic(errors.New("Unknown hash encoding type"))
	}
}

// 为对象的引用计数减一
// 当对象的引用计数降为0时，释放对象
func decrRefCount(robj *redisObject) {
	if robj.refcount <= 0 {
		panic(errors.New("decrRefCount against refcount <= 0"))
	}
	if robj.refcount == 1 {
		switch robj.rtype {
		case REDIS_STRING:
			freeStringObject(robj)
		case REDIS_LIST:
			freeListObject(robj)
		case REDIS_SET:
			freeSetObject(robj)
		case REDIS_ZSET:
			freeZsetObject(robj)
		case REDIS_HASH:
			freeHashObject(robj)
		default:
			panic(errors.New("Unknown object type"))
		}
		robj = nil
	} else {
		robj.refcount--
	}
}

func compareStringObjectsWithFlags(a *redisObject, b *redisObject, flags int) int {
	if a.rtype != REDIS_STRING || b.rtype != REDIS_STRING {
		panic(errors.New("type must redis string"))
	}
	if a == b {
		return 1
	}
	// todo
	return 0
}

func compareStringObjects(a *redisObject, b *redisObject) int {
	return compareStringObjectsWithFlags(a, b, REDIS_COMPARE_BINARY)
}
