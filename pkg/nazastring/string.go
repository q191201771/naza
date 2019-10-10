package nazastring

import "unsafe"

type sliceT struct {
	array unsafe.Pointer
	len   int
	cap   int
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}

func SliceByteToStringTmp(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToSliceByteTmp(s string) []byte {
	str := (*stringStruct)(unsafe.Pointer(&s))
	ret := sliceT{array: unsafe.Pointer(str.str), len: str.len, cap: str.len}
	return *(*[]byte)(unsafe.Pointer(&ret))
}
