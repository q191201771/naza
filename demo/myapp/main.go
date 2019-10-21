package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/q191201771/naza/pkg/bininfo"
)

func main() {
	v := flag.Bool("v", false, "show bin info")
	flag.Parse()
	if *v {
		_, _ = fmt.Fprint(os.Stderr, bininfo.StringifyMultiLine())
		os.Exit(1)
	}

	fmt.Println("my app running...")
	fmt.Println("bye...")
}
