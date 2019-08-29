// Package bininfo 将编译时的git版本号，时间，Go编译器信息打入程序中
package bininfo

import (
	"fmt"
	"runtime"
	"strings"
)

// 编译时通过如下方式传入编译时信息
//
// #GitCommitID=`git log --pretty=format:'%h' -n 1`
// GitCommitLog=`git log --pretty=oneline -n 1`
// GitStatus=`git status -s`
// BuildTime=`date +'%Y.%m.%d.%H%M%S'`
// BuildGoVersion=`go version`
//
// go build -ldflags " \
// -X 'github.com/q191201771/nezha/pkg/bininfo.GitCommitLog=${GitCommitLog}' \
// -X 'github.com/q191201771/nezha/pkg/bininfo.GitStatus=${GitStatus}' \
// -X 'github.com/q191201771/nezha/pkg/bininfo.BuildTime=${BuildTime}' \
// -X 'github.com/q191201771/nezha/pkg/bininfo.BuildGoVersion=${BuildGoVersion}' \
// "

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
