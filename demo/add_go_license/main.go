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

var license = `// Copyright %d, Chef.  All rights reserved.
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
	head := fmt.Sprintf(license, year, repo)
	nazalog.Debug(head)

	err := filebatch.Walk(dir, true, ".go", func(path string, info os.FileInfo, content []byte) []byte {
		return nil
	})
	nazalog.FatalIfErrorNotNil(err)
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
