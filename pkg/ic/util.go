package ic

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
	"sort"
)

type IDSlice []uint32

func (a IDSlice) Len() int           { return len(a) }
func (a IDSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a IDSlice) Less(i, j int) bool { return a[i] < a[j] }

func resetBuf(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	return b
}

func sortIDSlice(ids IDSlice) {
	sort.Sort(ids)
}

func zlibWrite(in []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(in)
	_ = w.Close()
	return b.Bytes()
}

func zlibRead(in []byte) (ret []byte) {
	b := bytes.NewReader(in)
	r, _ := zlib.NewReader(b)
	ret, _ = ioutil.ReadAll(r)
	return
}

//func isBufEmpty(b []byte) bool {
//	for i := 0; i < len(b); i++ {
//		if b[i] != 0 {
//			return false
//		}
//	}
//	return true
//}
//
//func dumpIDSlice(ids IDSlice, filename string) {
//	fp, _ := os.Create(filename)
//	for _, id := range ids {
//		_, _ = fp.WriteString(fmt.Sprintf("%d", id))
//		_, _ = fp.WriteString("\n")
//	}
//	_ = fp.Close()
//}
