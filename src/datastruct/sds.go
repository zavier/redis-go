package datastruct

// buf最大预分配长度
const SDS_MAX_PREALLOC = 1024 * 1024

// 直接使用sds为数据切片，sdshdr相关可以用切片的属性替代
// type slice struct {
//	array unsafe.Pointer
//	len   int
//	cap   int
//}
type sds []byte

// 返回已用空间大小
func sdsLen(s sds) int {
	return len(s)
}

// 返回可用空间大小
func sdsAvail(s sds) int {
	return cap(s) - len(s)
}

// 创建一个新的sds
func sdsNewLen(init interface{}, initlen int) sds {
	sds := make([]byte, initlen)
	if init != nil && initlen > 0 {
		switch init.(type) {
		case string:
			copy(sds, init.(string))
		case []byte:
			copy(sds, init.([]byte))
		}
	}
	return sds
}

// 创建一个空字符串
func sdsEmpty() sds {
	return sdsNewLen("", 0)
}

// 根据给定字符串创建sds
func sdsNew(init string) sds {
	return sdsNewLen(init, len(init))
}

// 创建sds副本
func sdsDup(s sds) sds {
	return sdsNewLen(string(s), len(s))
}

// 释放字符数组占用空间
func sdsFree(s sds) sds {
	return nil
}

// 将字符数组置空
func sdsClear(s sds) sds {
	return make([]byte, 0)
}

/*
重新为sds中数组分配长度，增加addlen个长度数组
如果加上后总长度大于 SDS_MAX_PREALLOC，则新长度为 (oldlen+addlen)*2 ，否则为 (oldlen+addlen)+SDS_MAX_PREALLOC
*/
func sdsMakeRoomFor(s sds, addlen int) sds {
	free := sdsAvail(s)
	if free > addlen {
		return s
	}
	newLen := sdsLen(s) + addlen
	if newLen < SDS_MAX_PREALLOC {
		newLen *= 2
	} else {
		newLen += SDS_MAX_PREALLOC
	}
	newSds := make([]byte, len(s), newLen)
	copy(newSds, s)
	return newSds
}

// 删除字符数组中的多余空间
func sdsRemoveFreeSpace(s sds) sds {
	newSds := make([]byte, len(s))
	copy(newSds, s)
	return newSds
}

// 返回给定sds分配的字节数
func sdsAllocSize(s sds) int {
	return cap(s)
}

// 扩大占用的空间，减少剩余空间（剩余空间足够的情况下）
// incr如果为负数则进行右截断
func sdsIncrLen(s sds, incr int) sds {
	if incr < 0 {
		newLen := sdsLen(s) + incr
		return s[0 : newLen-1]
	}

	availNum := sdsAvail(s)
	if incr < availNum {
		temp := make([]byte, incr)
		return append(s, temp...)
	}
	return s
}

// 将sds扩充至指定长度
func sdsGrowZero(s sds, len int) sds {
	curLen := sdsLen(s)
	if len <= curLen {
		return s
	}
	s = sdsMakeRoomFor(s, len-curLen)
	return s
}

// 将长度为len的字符串t追加到sds字符串末尾
func sdsCatLen(s sds, t string, len int) sds {
	s = sdsMakeRoomFor(s, len)
	s = append(s, []byte(t)...)
	return s
}
