package bufferpool

import (
	"bytes"
	"sync"
	"sync/atomic"
)

var minSize = 1024

type BufferPool struct {
	getCount uint32
	putCount uint32
	hitCount uint32
	mallocCount uint32

	m sync.Mutex
	// TODO chef: 这个map可以预申请，做成fixed size的
	sizeToList map[int]*[]*bytes.Buffer
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		sizeToList: make(map[int]*[]*bytes.Buffer),
	}
}

func (bp *BufferPool) Get(size int) *bytes.Buffer {
	atomic.AddUint32(&bp.getCount, 1)
	ss := up2power(size)
	if ss < minSize {
		ss = minSize
	}

	bp.m.Lock()
	l, ok := bp.sizeToList[ss]
	if !ok {
		bp.m.Unlock()
		return bp.newBuffer(ss)
	} else {
		if len(*l) == 0 {
			bp.m.Unlock()
			return bp.newBuffer(ss)
		}
		buf := (*l)[len(*l)-1]
		*l = (*l)[:len(*l)-1]
		bp.m.Unlock()
		buf.Reset()
		atomic.AddUint32(&bp.hitCount, 1)
		return buf
	}
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
	atomic.AddUint32(&bp.putCount, 1)
	size := down2power(buf.Cap())
	if size < minSize {
		size = minSize
	}

	bp.m.Lock()
	l, ok := bp.sizeToList[size]
	if !ok {
		l = new([]*bytes.Buffer)
		*l = append(*l, buf)
		// TODO
		bp.sizeToList[size] = l
	} else {
		*l = append(*l, buf)
	}
	bp.m.Unlock()
}

func (bp *BufferPool) newBuffer(n int) *bytes.Buffer {
	var buf bytes.Buffer
	buf.Grow(n)
	atomic.AddUint32(&bp.mallocCount, 1)
	return &buf
}

// @return 范围为 [2, 4, 8, 16, ..., 1073741824]，如果大于等于1073741824，则直接返回n
func up2power(n int) int {
	if n >= 1073741824 {
		return n
	}

	var i uint32
	for ; n > (2 << i); i++ {
	}
	return 2 << i
}

// @return 范围为 [2, 4, 8, 16, ..., 1073741824]
func down2power(n int) int {
	if n < 2 {
		return 2
	} else if n >= 1073741824 {
		return 1073741824
	}

	var i uint32
	for {
		nn := 2 << i
		if n > nn {
			i++
		} else if n == nn {
			return n
		} else if n < nn {
			return 2 << (i - 1)
		}
	}
}
