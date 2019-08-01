package datastruct

import (
	"strconv"
	"strings"
)

// buf最大预分配长度
const SDS_MAX_PREALLOC = 1024 * 1024

// *****************直接使用sds为数据切片，sdshdr相关可以用切片的属性替代*****************
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
	s = nil
	return s
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

// 将指定字符串添加到sds末尾
func sdsCat(s sds, t string) sds {
	return sdsCatLen(s, t, len(s))
}

// 将t添加到sds末尾
func sdsCatSds(s sds, t sds) {
	sdsCatLen(s, string(t), sdsLen(t))
}

// 将字符串t的前length个字符复制到sds s中,覆盖原有字符串
func sdsCpyLen(s sds, t string, length int) sds {
	totlen := cap(s)
	if totlen < length {
		s = sdsMakeRoomFor(s, length-len(s))
	}
	copy(s, s[0:length])
	return s
}

// 将字符串t复制到sds s中,覆盖原有字符串
func sdsCpy(s sds, t string) sds {
	return sdsCpyLen(s, t, len(t))
}

const SDS_LLSTR_SIZE = 21

// 原 sdsull2str(char *s, u long long v)
func sdsInt2Str(v uint) (string, int) {
	str := strconv.Itoa(int(v))
	return str, len(str)
}

// 根据输入的数字创建一个SDS, 原 sdsfromlonglong(long long value)
func sdsFromInt(value int) sds {
	str, length := sdsInt2Str(uint(value))
	buf := make([]byte, SDS_LLSTR_SIZE)
	copy(buf, str)
	return sdsNewLen(buf, length)
}

// 打印函数
func sdsCatVPrinf(s sds, fmt string) {
	//todo
}

// 打印函数
func sdsCatPrintf(s sds, fmt string, s1 ...interface{}) {
	//todo
}

// 对sds左右两端进行裁剪，清楚两端的所有cset中出现的字符
// s = sdsnew("AA...AA.a.aa.aHelloWorld     :::");
// s = sdstrim(s,"Aa. :");
// printf("%s\n", s); = HelloWorld
func sdsTrim(s sds, cset string) sds {
	newStr := strings.Trim(string(s), cset)
	return sdsNew(newStr)
}

// 裁剪sds, 索引可以为负数
func sdsRange(s sds, start int, end int) sds {
	str := string(s)
	if start < 0 {
		start = len(str) + start
		if start < 0 {
			start = 0
		}
	}
	if end < 0 {
		end = len(str) + end
		if end < 0 {
			end = 0
		}
	}
	newStr := str[start:end]
	len := end - start + 1
	newSds := make([]byte, len)
	copy(newSds, newStr)
	return newSds
}

// 将sds字符串中的所有字符转小写
func sdsToLower(s sds) sds {
	str := string(s)
	newStr := strings.ToLower(str)
	return sdsNew(newStr)
}

// 转大写
func sdsToUpper(s sds) sds {
	str := string(s)
	newStr := strings.ToUpper(str)
	return sdsNew(newStr)
}

// 比较s1与s2
func sdsCmp(s1 sds, s2 sds) int {
	return strings.Compare(string(s1), string(s2))
}

// 使用sep对s进行分割
func sdsSplitLen(s string, sep string) ([]sds, int) {
	split := strings.Split(s, sep)
	num := len(split)
	sdss := make([]sds, 0)
	for _, str := range split {
		sdss = append(sdss, sdsNew(str))
	}
	return sdss, num
}

// 释放count个sds
func sdsFreeSplitRes(tokens []sds, count int) []sds {
	if count > len(tokens) {
		return make([]sds, 0)
	}
	remain := len(tokens) - count
	newSds := tokens[count:]
	res := make([]sds, remain)
	copy(res, newSds)
	return res
}

// todo
func sdsCatRepr(s sds, p string) {

}

// 是否是16进制符号中的一个
func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}

// 16进制数转10进制
func hexDigitToInt(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c) - 48
	}
	if c >= 'a' && c <= 'f' {
		return int(c) - 97 + 10
	}
	return 0
}

// todo
func sdsSplitArgs(line string, argc int) {

}

// 将sds中 from中的字符串替换为to字符串
func sdsMapChars(s sds, from string, to string) sds {
	str := string(s)
	replaceRes := strings.Replace(str, from, to, -1)
	return sdsNew(replaceRes)
}

// join
func sdsJoin(str []string, sep string) sds {
	join := strings.Join(str, sep)
	return sdsNew(join)
}
