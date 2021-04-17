// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package filesystemlayer_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/assert"

	"github.com/q191201771/naza/pkg/filesystemlayer"
)

func TestCase1(t *testing.T) {
	fslCtx := filesystemlayer.FSLFactory(filesystemlayer.FSLTypeMemory)

	var wg sync.WaitGroup
	wg.Add(16)

	for i := 0; i < 16; i++ {
		go func(ii int) {
			dir := fmt.Sprintf("/tmp/lal/hls/test%d", ii)
			err := fslCtx.MkdirAll(dir, 0777)
			assert.Equal(t, nil, err)

			for j := 0; j < 32; j++ {
				filename := fmt.Sprintf("/tmp/lal/hls/test%d/%d.ts", ii, j)
				nazalog.Infof("%d %d %s", ii, j, filename)
				fp, err := fslCtx.Create(filename)
				assert.Equal(t, nil, err)

				n, err := fp.Write([]byte("hello"))
				assert.Equal(t, nil, err)
				assert.Equal(t, 5, n)

				n, err = fp.Write([]byte("world"))
				assert.Equal(t, nil, err)
				assert.Equal(t, 5, n)

				err = fp.Close()
				assert.Equal(t, nil, err)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	// 正常读
	b, err := fslCtx.ReadFile("/tmp/lal/hls/test1/1.ts")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("helloworld"), b)

	// 删文件
	err = fslCtx.Remove("/tmp/lal/hls/test1/1.ts")
	assert.Equal(t, nil, err)
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test1/1.ts")
	assert.Equal(t, filesystemlayer.ErrNotFound, err)
	assert.Equal(t, nil, b)

	// 正常读
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test2/2.ts")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("helloworld"), b)

	// 文件重命名
	err = fslCtx.Rename("/tmp/lal/hls/test2/2.ts", "/tmp/lal/hls/test2/new2.ts")
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test2/2.ts")
	assert.Equal(t, filesystemlayer.ErrNotFound, err)
	assert.Equal(t, nil, b)
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test2/new2.ts")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("helloworld"), b)

	// 删文件夹
	err = fslCtx.RemoveAll("/tmp/lal/hls/test1")
	assert.Equal(t, nil, err)
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test1/1.ts")
	assert.Equal(t, filesystemlayer.ErrNotFound, err)
	assert.Equal(t, nil, b)

	// 创建已经存在的文件
	fp, err := fslCtx.Create("/tmp/lal/hls/test3/3.ts")
	assert.Equal(t, nil, err)
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test3/3.ts")
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, b)
	n, err := fp.Write([]byte("asd"))
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, n)
	b, err = fslCtx.ReadFile("/tmp/lal/hls/test3/3.ts")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("asd"), b)

	// 删除不存在的文件
	err = fslCtx.Remove("/tmp/lal/hls/test1/1.ts")
	assert.Equal(t, filesystemlayer.ErrNotFound, err)

	// 重命名不存在的文件
	err = fslCtx.Rename("/tmp/lal/hls/test1/1.ts", "/tmp/lal/hls/test1/new1.ts")
	assert.Equal(t, filesystemlayer.ErrNotFound, err)
}
