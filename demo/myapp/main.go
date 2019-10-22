// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

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
