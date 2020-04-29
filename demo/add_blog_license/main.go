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
	"os"

	"github.com/q191201771/naza/pkg/filebatch"
	"github.com/q191201771/naza/pkg/nazalog"
)

//var licenseTmpl = `
//> **原文链接：** [https://pengrl.com/p/%s/](https://pengrl.com/p/%s/)
//> **原文出处：** [yoko blog](https://pengrl.com) (https://pengrl.com)
//> **原文作者：** [yoko](https://github.com/q191201771) (https://github.com/q191201771)
//> **版权声明：** 本文欢迎任何形式转载，转载时完整保留本声明信息（包含原文链接、原文出处、原文作者、版权声明）即可。本文后续所有修改都会第一时间在原始地址更新。
//
//![fccxy](https://pengrl.com/images/fccxy_qccode_and_sys.jpg)`

var licenseTmpl = `
本文完，作者[yoko](https://github.com/q191201771)，尊重劳动人民成果，转载请注明原文出处： [https://pengrl.com/p/%s/](https://pengrl.com/p/%s/)`

func main() {
	dir := parseFlag()

	//linesOfLicense := strings.Split(licenseTmpl, "\n")
	//lastLineOfLicense := linesOfLicense[len(linesOfLicense)-1]

	var (
		skipCount int
		modCount  int
	)
	err2 := filebatch.Walk(dir, true, ".md", func(path string, info os.FileInfo, content []byte, err error) []byte {
		if err != nil {
			nazalog.Warnf("read file failed. file=%s, err=%+v", path, err)
			return nil
		}
		lines := bytes.Split(content, []byte{'\n'})

		//if bytes.Index(lines[len(lines)-1], []byte("本文完")) != -1 ||
		//	res, err := filebatch.DeleteLines(content, filebatch.LineRange{From: -1, To: -1})
		//	nazalog.Debugf("%s -2", info.Name())
		//	nazalog.FatalIfErrorNotNil(err)
		//	return res
		//}
		//nazalog.Warnf("%s", info.Name())
		//return content

		// 已添加过声明，不用再添加了
		if bytes.Index(lines[len(lines)-1], []byte("本文完，作者")) != -1 ||
			bytes.Index(lines[len(lines)-2], []byte("本文完，作者")) != -1 {
			nazalog.Debug(info.Name())
			skipCount++
			return nil
		}

		// 获取该文章的url的地址
		var abbrlink string
		for _, line := range lines {
			if bytes.Index(line, []byte("abbrlink")) != -1 {
				abbrlink = string(bytes.TrimSpace(bytes.Split(line, []byte{':'})[1]))
				nazalog.Debug(abbrlink)
				break
			}
		}
		if abbrlink == "" {
			nazalog.Errorf("abbrlink not exist. path=%s", path)
			os.Exit(1)
		}

		// 构造好license信息，并添加在文章末尾
		modCount++
		license := fmt.Sprintf(licenseTmpl, abbrlink, abbrlink)
		return filebatch.AddTailContent(content, []byte(license))
	})
	nazalog.Assert(nil, err2)
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
