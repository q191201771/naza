package bufferpool

import (
	"bytes"
	"github.com/q191201771/naza/pkg/assert"
	"math/rand"
	"testing"
	"time"
)

var bp *BufferPool
var count int

func TestBufferPool(t *testing.T) {
	bp := NewBufferPool()
	buf := &bytes.Buffer{}
	bp.Put(buf)
	buf = bp.Get(4096)
	buf.Grow(4096)
	bp.Put(buf)
	buf = bp.Get(4096)
	bp.Put(buf)
}

func size() int {
	//return 1024

	//ss := []int{1000, 2000, 5000}
	////ss := []int{128, 1024, 4096, 16384}
	//count++
	//return ss[count % 3]

	return random(0, 128 * 1024)
}

func random(l, r int) int {
	return l + (rand.Int() % (r - l))
}

func origin() {
	var buf bytes.Buffer
	size := size()
	buf.Grow(size)
}

func bufferPool() {
	size := size()
	buf := bp.Get(size)
	buf.Grow(size)
	bp.Put(buf)
}

func BenchmarkOrigin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origin()
	}
}

func BenchmarkBufferPool(b *testing.B) {
	bp = NewBufferPool()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bufferPool()
	}
}

func BenchmarkOriginParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				origin()
			}
		})
	}
}

func BenchmarkBufferPoolParallel(b *testing.B) {
	bp = NewBufferPool()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bufferPool()
			}
		})
	}
}

func TestUp2power(t *testing.T) {
	assert.Equal(t, 2, up2power(0))
	assert.Equal(t, 2, up2power(1))
	assert.Equal(t, 2, up2power(2))
	assert.Equal(t, 4, up2power(3))
	assert.Equal(t, 4, up2power(4))
	assert.Equal(t, 8, up2power(5))
	assert.Equal(t, 8, up2power(6))
	assert.Equal(t, 8, up2power(7))
	assert.Equal(t, 8, up2power(8))
	assert.Equal(t, 16, up2power(9))
	assert.Equal(t, 1024, up2power(1023))
	assert.Equal(t, 1024, up2power(1024))
	assert.Equal(t, 2048, up2power(1025))
	assert.Equal(t, 1073741824, up2power(1073741824-1))
	assert.Equal(t, 1073741824, up2power(1073741824))
	assert.Equal(t, 1073741824+1, up2power(1073741824+1))
	assert.Equal(t, 2047483647-1, up2power(2047483647-1))
	assert.Equal(t, 2047483647, up2power(2047483647))
}

func TestDown2power(t *testing.T) {
	assert.Equal(t, 2, down2power(0))
	assert.Equal(t, 2, down2power(1))
	assert.Equal(t, 2, down2power(2))
	assert.Equal(t, 2, down2power(3))
	assert.Equal(t, 4, down2power(4))
	assert.Equal(t, 4, down2power(5))
	assert.Equal(t, 4, down2power(6))
	assert.Equal(t, 4, down2power(7))
	assert.Equal(t, 8, down2power(8))
	assert.Equal(t, 8, down2power(9))
	assert.Equal(t, 512, down2power(1023))
	assert.Equal(t, 1024, down2power(1024))
	assert.Equal(t, 1024, down2power(1025))
	assert.Equal(t, 1073741824 >> 1, down2power(1073741824-1))
	assert.Equal(t, 1073741824, down2power(1073741824))
	assert.Equal(t, 1073741824, down2power(1073741824+1))
	assert.Equal(t, 1073741824, down2power(2047483647-1))
	assert.Equal(t, 1073741824, down2power(2047483647))
}

func init() {
	rand.Seed(time.Now().Unix())
}
