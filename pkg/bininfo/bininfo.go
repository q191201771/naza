// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

// Package bininfo 将编译时的 git commit 日志，时间，Go 编译器信息打入程序中
package bininfo

import (
	"fmt"
	"runtime"
	"strings"
)

// 编译时通过如下方式传入编译时信息
//
// GitCommitLog=`git log --pretty=oneline -n 1`
// # 将 log 原始字符串中的单引号替换成双引号
// GitCommitLog=${GitCommitLog//\'/\"}
//
// GitStatus=`git status -s`
// BuildTime=`date +'%Y.%m.%d.%H%M%S'`
// BuildGoVersion=`go version`
//
// LDFlags=" \
//     -X 'github.com/q191201771/naza/pkg/bininfo.GitCommitLog=${GitCommitLog}' \
//     -X 'github.com/q191201771/naza/pkg/bininfo.GitStatus=${GitStatus}' \
//     -X 'github.com/q191201771/naza/pkg/bininfo.BuildTime=${BuildTime}' \
//     -X 'github.com/q191201771/naza/pkg/bininfo.BuildGoVersion=${BuildGoVersion}' \
// "
//
// go build -ldflags "$LDFlags"

var (
	GitCommitLog   = "unknown"
	GitStatus      = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

func StringifySingleLine() string {
	return fmt.Sprintf("GitCommitLog=%s. GitStatus=%s. BuildTime=%s. GoVersion=%s. runtime=%s/%s.",
		GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func StringifyMultiLine() string {
	return fmt.Sprintf("GitCommitLog=%s\nGitStatus=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		GitCommitLog, GitStatus, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func beauty() {
	if GitStatus == "" {
		GitStatus = "cleanly"
	} else {
		GitStatus = strings.Replace(strings.Replace(GitStatus, "\r\n", " |", -1), "\n", " |", -1)
	}
}

func init() {
	beauty()
}
