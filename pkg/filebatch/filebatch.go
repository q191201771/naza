package filebatch

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// @param path 带路径的文件名
// @param info 文件的 os.FileInfo 信息
// @param content 文件内容
// @return 返回nil或者content原始内容，则不修改文件内容，返回其他内容，则会覆盖重写文件
type WalkFunc func(path string, info os.FileInfo, content []byte) []byte

// 遍历访问指定文件夹下的文件
// @param root 需要遍历访问的文件夹
// @param recursive 是否递归访问子文件夹
// @param suffix 指定文件名后缀进行过滤，如果为""，则不过滤
func Walk(root string, recursive bool, suffix string, walkFn WalkFunc) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !recursive && info.IsDir() && path != root {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if suffix != "" && !strings.HasSuffix(info.Name(), suffix) {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		newContent := walkFn(path, info, content)
		if newContent != nil && bytes.Compare(content, newContent) != 0 {
			if err = ioutil.WriteFile(path, newContent, 0755); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// 文件尾部添加内容
func AddTailContent(content []byte, tail []byte) []byte {
	if !bytes.HasSuffix(content, []byte{'\n'}) {
		content = append(content, '\n')
	}
	return append(content, tail...)
}

func AddHeadContent(content []byte, head []byte) []byte {
	if !bytes.HasSuffix(head, []byte{'\n'}) {
		head = append(head, '\n')
	}
	return append(head, content...)
}
