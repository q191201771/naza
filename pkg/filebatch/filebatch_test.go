// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package filebatch

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

var filenameToContent map[string][]byte

var head = `// Copyright %s, Chef.  All rights reserved.
// https://%s
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)`

var tail = `
> author: xxx
> link: xxx
> license: xxx
`

// /<root>/
//     |-- /dir1/
//     |-- /dir2/
//         |-- file5
//         |-- file6
//         |-- file7.txt
//         |-- file8.txt
//     |-- file1
//     |-- file2
//     |-- file3.txt
//     |-- file4.txt
func prepareTestFile() (string, error) {
	filenameToContent = make(map[string][]byte)

	root, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	if root[len(root)-1] != '/' {
		root = root + "/"
	}
	nazalog.Debugf(root)

	if err = os.Mkdir(filepath.Join(root, "dir1"), 0755); err != nil {
		return "", err
	}
	if err = os.Mkdir(filepath.Join(root, "dir2"), 0755); err != nil {
		return "", err
	}

	filenameToContent[root+"file1"] = []byte("hello")
	filenameToContent[root+"file2"] = []byte("hello")
	filenameToContent[root+"file3.txt"] = []byte("hello")
	filenameToContent[root+"file4.txt"] = []byte("hello")
	filenameToContent[root+"dir2/file5"] = []byte("hello")
	filenameToContent[root+"dir2/file6"] = []byte("hello")
	filenameToContent[root+"dir2/file7.txt"] = []byte("hello")
	filenameToContent[root+"dir2/file8.txt"] = []byte("hello")

	for k, v := range filenameToContent {
		if err = ioutil.WriteFile(k, v, 0755); err != nil {
			return "", err
		}
	}

	return root, nil
}

func testWalk(t *testing.T, recursive bool, suffix string) {
	root, err := prepareTestFile()
	assert.Equal(t, nil, err)
	defer os.RemoveAll(root)

	err2 := Walk(root, recursive, suffix, func(path string, info os.FileInfo, content []byte, err error) []byte {
		nazalog.Debugf("%+v %+v %s", path, info.Name(), string(content))

		v := filenameToContent[path]
		assert.Equal(t, v, content)
		delete(filenameToContent, path)

		return content
	})
	assert.Equal(t, nil, err2)
}

func TestWalk(t *testing.T) {
	testWalk(t, true, "")
	assert.Equal(t, 0, len(filenameToContent))

	testWalk(t, false, "")
	assert.Equal(t, 4, len(filenameToContent))

	testWalk(t, true, ".txt")
	assert.Equal(t, 4, len(filenameToContent))

	testWalk(t, false, ".txt")
	assert.Equal(t, 6, len(filenameToContent))

	testWalk(t, false, ".notexist")
	assert.Equal(t, 8, len(filenameToContent))
}

func TestAddContent(t *testing.T) {
	root, err := prepareTestFile()
	assert.Equal(t, nil, err)
	defer os.RemoveAll(root)

	err2 := Walk(root, true, ".txt", func(path string, info os.FileInfo, content []byte, err error) []byte {
		lines := bytes.Split(content, []byte{'\n'})
		nazalog.Debugf("%+v %d", path, len(lines))

		v := filenameToContent[path]
		assert.Equal(t, v, content)
		delete(filenameToContent, path)

		return AddHeadContent(AddTailContent(content, []byte(tail)), []byte(head))
	})
	assert.Equal(t, nil, err2)

	err2 = Walk(root, true, "", func(path string, info os.FileInfo, content []byte, err error) []byte {
		nazalog.Debugf("%+v %+v %s", path, info.Name(), string(content))
		return nil
	})
	assert.Equal(t, nil, err2)
}

func TestDeleteLines(t *testing.T) {
	origin := `111
222
333
444
555`

	content := []byte(origin)
	lines := bytes.Split(content, []byte{'\n'})
	assert.Equal(t, 5, len(lines))

	var (
		res []byte
		err error
	)

	// 常规操作
	res, err = DeleteLines(content, LineRange{From: 1, To: 1})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`222
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: -5, To: -5})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`222
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: 2, To: 2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: -4, To: -4})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
333
444
555`), res)

	res, err = DeleteLines(content, LineRange{From: 4, To: 4})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
555`), res)

	res, err = DeleteLines(content, LineRange{From: -2, To: -2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
555`), res)

	res, err = DeleteLines(content, LineRange{From: 5, To: 5})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
444`), res)

	res, err = DeleteLines(content, LineRange{From: -1, To: -1})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222
333
444`), res)

	res, err = DeleteLines(content, LineRange{From: 1, To: 3})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`444
555`), res)

	res, err = DeleteLines(content, LineRange{From: -5, To: -3})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`444
555`), res)

	res, err = DeleteLines(content, LineRange{From: 3, To: 5})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222`), res)

	res, err = DeleteLines(content, LineRange{From: -3, To: -1})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
222`), res)

	res, err = DeleteLines(content, LineRange{From: 2, To: 4})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
555`), res)

	res, err = DeleteLines(content, LineRange{From: -4, To: -2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
555`), res)

	// 非常规操作
	res, err = DeleteLines(content, LineRange{From: 4, To: 2})
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte(`111
555`), res)

	res, err = DeleteLines(content, LineRange{From: 0, To: 1})
	assert.Equal(t, ErrLineRange, err)

	res, err = DeleteLines(content, LineRange{From: 1, To: 0})
	assert.Equal(t, ErrLineRange, err)

	res, err = DeleteLines(content, LineRange{From: 10, To: 20})
	assert.Equal(t, ErrLineRange, err)
}
