// Package bininfo 将编译时的git版本号，时间，Go编译器信息打入程序中
package bininfo

import (
	"fmt"
	"runtime"
)

// 编译时通过如下方式传入编译时信息
// go build -ldflags " \
//   -X 'github.com/q191201771/nezha/pkg/bininfo/bininfo.GitCommitID=`git log --pretty=format:'%h' -n 1`' \
//   -X 'github.com/q191201771/nezha/pkg/bininfo/bininfo.BuildTime=`date +'%Y.%m.%d.%H%M%S'`' \
//   -X 'github.com/q191201771/nezha/pkg/bininfo/bininfo.BuildGoVersion=`go version`' \
// "

var (
	GitCommitID    string
	BuildTime      string
	BuildGoVersion string
)

func StringifySingleLine() string {
	return fmt.Sprintf("GitCommitID=%s. BuildTime=%s. GoVersion=%s. runtime=%s/%s.",
		GitCommitID, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func StringifyMultiLine() string {
	return fmt.Sprintf("GitCommitID=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s.",
		GitCommitID, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func init() {
	if GitCommitID == "" {
		GitCommitID = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if BuildGoVersion == "" {
		BuildGoVersion = "unknown"
	}
}
