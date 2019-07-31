package datastruct

import (
	"testing"
)

func TestSdsNew(t *testing.T) {
	ori := sdsNew("abcdef")
	if len(ori) != 6 || cap(ori) != 6 {
		t.Errorf("error, len %d, cap %d", len(ori), cap(ori))
	}
}

func TestSdsMakeRoomFor(t *testing.T) {
	ori := sdsNew("abcdefg")
	newSds := sdsMakeRoomFor(ori, 10)
	newLen := (len(ori) + 10) * 2
	if len(newSds) != 7 || cap(newSds) != newLen {
		t.Errorf("error, len %d ,cap %d", len(newSds), cap(newSds))
	}
}

func TestHexDigitToInt(t *testing.T) {
	i := hexDigitToInt('1')
	if i != 1 {
		t.Error("error != 1")
	}
	toInt := hexDigitToInt('b')
	if toInt != 11 {
		t.Error("err != 11")
	}
}
