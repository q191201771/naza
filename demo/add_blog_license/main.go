package main

import (
	"bytes"
	"flag"
	"github.com/q191201771/naza/pkg/filebatch"
	"github.com/q191201771/naza/pkg/nazalog"
	"os"
)

func main() {
	dir := parseFlag()
	err := filebatch.Walk(dir, true, ".md", func(path string, info os.FileInfo, content []byte) []byte {
		lines := bytes.Split(content, []byte{'\n'})
		nazalog.Debug(path, len(lines))
		return nil
	})
	nazalog.FatalIfErrorNotNil(err)
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
