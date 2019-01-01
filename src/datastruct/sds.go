package datastruct

type sds string

const SDS_MAX_PREALLOC = 1024 * 1024

/*
go 语言无char类型且提供了string，此处直接使用string类型
*/
type sdshdr8 struct {
	len   uint8  // used
	alloc uint8  // excluding the header
	flags uint8  // 3 1sb of type, 5 unused bits
	buf   string // []byte
}

/*
创建一个新的sds（string）
*/
func (self *sds) sdsNewLen(init string, initlen uint64) sds {
	// 原代码为根据长度分配不同的sdshdr8, sdshdr16等
	return sds(init)
}

func (self *sds) sdsEmpty() sds {
	return self.sdsNewLen("", 0)
}

func (self *sds) sdsNew(init string) sds {
	return self.sdsNewLen(init, uint64(len(init)))
}

func (self *sds) sdsUp(s sds) sds {
	return self.sdsNewLen(string(s), uint64(len(s)))
}

/*
释放字符数组占用空间
*/
func (self *sds) sdsFree(s sds) {
	s = ""
}

/*
更新sds的长度为len(s)
*/
func (self *sds) sdsUpdateLen(s sds) {
}

/*
将字符数组置空
*/
func (self *sds) sdsClear(s sds) {
	s = ""
}

/*
重新为sds中数组分配长度，增加addlen个长度数组
如果加上后总长度大于 SDS_MAX_PREALLOC，则新长度为 (oldlen+addlen)*2 ，否则为 (oldlen+addlen)+SDS_MAX_PREALLOC
根据新长度判断sds类型，
*/
func (self *sds) sdsMakeRoomFor(s sds, addlen uint64) sds {
	return s
}

/*
删除字符数组中的多余空间
*/
func (self *sds) sdsRemoveFreeSpace(s sds) sds {
	return s
}

func (self *sds) sdaAllocSize(s sds) uint64 {
	return uint64(len(string(s)))
}
