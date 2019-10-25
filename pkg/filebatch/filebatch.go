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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// @param path 带路径的文件名
// @param info 文件的 os.FileInfo 信息
// @param content 文件内容
// @return 返回nil或者content原始内容，则不修改文件内容，返回其他内容，则会覆盖重写文件
type WalkFunc func(path string, info os.FileInfo, content []byte, err error) []byte

// 遍历访问指定文件夹下的文件
// @param root 需要遍历访问的文件夹
// @param recursive 是否递归访问子文件夹
// @param suffix 指定文件名后缀进行过滤，如果为""，则不过滤
func Walk(root string, recursive bool, suffix string, walkFn WalkFunc) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			walkFn(path, info, nil, err)
			return nil
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
			walkFn(path, info, content, err)
			return nil
		}
		newContent := walkFn(path, info, content, nil)
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

// 文件头部添加内容
func AddHeadContent(content []byte, head []byte) []byte {
	if !bytes.HasSuffix(head, []byte{'\n'}) {
		head = append(head, '\n')
	}
	return append(head, content...)
}

// 行号范围
// 1表示首行，-1表示最后一行
type LineRange struct {
	From int
	To   int
}

var ErrLineRange = errors.New("naza.filebatch: line range error")

func calcLineRange(len int, lr LineRange) (LineRange, error) {
	// 换算成从0开始的下标
	if lr.From < 0 {
		lr.From = len + lr.From
	} else if lr.From > 0 {
		lr.From = lr.From - 1
	} else {
		return lr, ErrLineRange
	}
	if lr.To < 0 {
		lr.To = len + lr.To
	} else if lr.To > 0 {
		lr.To = lr.To - 1
	} else {
		return lr, ErrLineRange
	}

	// 排序交换
	if lr.From > lr.To {
		lr.From, lr.To = lr.To, lr.From
	}

	if lr.From < 0 || lr.From >= len || lr.To < 0 || lr.To >= len {
		return lr, ErrLineRange
	}

	return lr, nil
}

func DeleteLines(content []byte, lr LineRange) ([]byte, error) {
	lines := bytes.Split(content, []byte{'\n'})
	length := len(lines)
	nlr, err := calcLineRange(length, lr)
	if err != nil {
		return content, err
	}
	var nlines [][]byte
	if nlr.From > 0 {
		nlines = append(nlines, lines[:nlr.From]...)
	}
	if nlr.To < length-1 {
		nlines = append(nlines, lines[nlr.To+1:]...)
	}
	return bytes.Join(nlines, []byte{'\n'}), nil
}
