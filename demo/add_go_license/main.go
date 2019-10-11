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
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

var licenseTmpl = `// Copyright %d, Chef.  All rights reserved.
// https://%s
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

`

func main() {
	dir := parseFlag()

	year := time.Now().Year()
	repo := achieveRepo(dir)
	license := fmt.Sprintf(licenseTmpl, year, repo)
	nazalog.Debug(license)

	var (
		skipCount int
		modCount  int
	)
	err := filebatch.Walk(dir, true, ".go", func(path string, info os.FileInfo, content []byte) []byte {
		lines := bytes.Split(content, []byte{'\n'})
		if bytes.Index(lines[0], []byte("Copyright")) != -1 {
			skipCount++
			//nc, _ := filebatch.DeleteLines(content, filebatch.LineRange{From:1, To:7})
			//return nc
			return nil
		}

		modCount++
		return filebatch.AddHeadContent(content, []byte(license))
	})
	nazalog.FatalIfErrorNotNil(err)
	nazalog.Infof("count. mod=%d, skip=%d", modCount, skipCount)
}

func achieveRepo(root string) string {
	content, err := ioutil.ReadFile(filepath.Join(root, "go.mod"))
	nazalog.FatalIfErrorNotNil(err)
	lines := bytes.Split(content, []byte{'\n'})
	repo := bytes.TrimPrefix(lines[0], []byte("module "))
	return string(bytes.TrimSpace(repo))
}

func parseFlag() string {
	dir := flag.String("d", "", "dir of repo")
	flag.Parse()
	if *dir == "" {
		flag.Usage()
		os.Exit(1)
	}
	return *dir
}
