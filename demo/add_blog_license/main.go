// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/q191201771/naza/pkg/filebatch"
	"github.com/q191201771/naza/pkg/nazalog"
	"os"
)

var licenseTmpl = `
> **本文原始地址：** [https://pengrl.com/p/%s/](https://pengrl.com/p/%s/)  
> **声明：** 本文后续所有修改都会第一时间在原始地址更新。本文欢迎任何形式转载，转载时注明原始出处即可。`

func main() {
	dir := parseFlag()

	//linesOfLicense := strings.Split(licenseTmpl, "\n")
	//lastLineOfLicense := linesOfLicense[len(linesOfLicense)-1]

	var (
		skipCount int
		modCount  int
	)
	err := filebatch.Walk(dir, true, ".md", func(path string, info os.FileInfo, content []byte) []byte {
		lines := bytes.Split(content, []byte{'\n'})
		if bytes.Index(lines[len(lines)-1], []byte("声明")) != -1 ||
			bytes.Index(lines[len(lines)-2], []byte("声明")) != -1 {
			skipCount++
			return nil

		}
		var abbrlink string
		for _, line := range lines {
			if bytes.Index(line, []byte("abbrlink")) != -1 {
				abbrlink = string(bytes.TrimSpace(bytes.Split(line, []byte{':'})[1]))
				nazalog.Debug(abbrlink)
				break
			}
		}

		modCount++
		license := fmt.Sprintf(licenseTmpl, abbrlink, abbrlink)
		return filebatch.AddTailContent(content, []byte(license))
	})
	nazalog.FatalIfErrorNotNil(err)
	nazalog.Infof("count. mod=%d, skip=%d", modCount, skipCount)
}

func parseFlag() string {
	dir := flag.String("d", "", "dir of posts")
	flag.Parse()
	if *dir == "" {
		flag.Usage()
		os.Exit(1)
	}
	return *dir
}
