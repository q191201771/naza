package filebatch

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/nazalog"
)

var filenameToContent map[string][]byte

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

	assert.Equal(t, nil, err)
	err = Walk(root, recursive, suffix, func(path string, info os.FileInfo, content []byte) []byte {
		nazalog.Debugf("%+v %+v %s", path, info.Name(), string(content))

		v := filenameToContent[path]
		assert.Equal(t, v, content)
		delete(filenameToContent, path)

		return content
	})
	assert.Equal(t, nil, err)
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
head := `// Copyright %s, Chef.  All rights reserved.
// https://%s
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)
`

	tail := `
> author: xxx
> link: xxx
> license: xxx
`

	root, err := prepareTestFile()
	assert.Equal(t, nil, err)
	defer os.RemoveAll(root)

	assert.Equal(t, nil, err)
	err = Walk(root, true, ".txt", func(path string, info os.FileInfo, content []byte) []byte {
		v := filenameToContent[path]
		assert.Equal(t, v, content)
		delete(filenameToContent, path)

		return AddHeadContent(AddTailContent(content, []byte(tail)), []byte(head))
	})
	assert.Equal(t, nil, err)

	err = Walk(root, true, "", func(path string, info os.FileInfo, content []byte) []byte {
		nazalog.Debugf("%+v %+v %s", path, info.Name(), string(content))
		return nil
	})
}
