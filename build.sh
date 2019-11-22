#!/usr/bin/env bash

set -x

# 获取源码最近一次 git commit log，包含 commit sha 值，以及 commit message
GitCommitLog=`git log --pretty=oneline -n 1`
# 将 log 原始字符串中的单引号替换成双引号
GitCommitLog=${GitCommitLog//\'/\"}
# 检查源码在git commit 基础上，是否有本地修改，且未提交的内容
GitStatus=`git status -s`
# 获取当前时间
BuildTime=`date +'%Y.%m.%d.%H%M%S'`
# 获取 Go 的版本
BuildGoVersion=`go version`

# 将以上变量序列化至 LDFlags 变量中
LDFlags=" \
    -X 'github.com/q191201771/naza/pkg/bininfo.GitCommitLog=${GitCommitLog}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.GitStatus=${GitStatus}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.BuildTime=${BuildTime}' \
    -X 'github.com/q191201771/naza/pkg/bininfo.BuildGoVersion=${BuildGoVersion}' \
"

ROOT_DIR=`pwd`

# 如果可执行程序输出目录不存在，则创建
if [ ! -d ${ROOT_DIR}/bin ]; then
  mkdir bin
fi

# 编译多个可执行程序
cd ${ROOT_DIR}/demo/add_blog_license && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/add_blog_license &&
cd ${ROOT_DIR}/demo/add_go_license && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/add_go_license &&
cd ${ROOT_DIR}/demo/taskpool && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/taskpool &&
cd ${ROOT_DIR}/demo/slicebytepool && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/slicebytepool &&
cd ${ROOT_DIR}/demo/myapp && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/myapp &&
#cd ${ROOT_DIR}/demo/time && go build -ldflags "$LDFlags" -o ${ROOT_DIR}/bin/time &&
ls -lrt ${ROOT_DIR}/bin &&
cd ${ROOT_DIR} && ./bin/myapp -v &&
echo 'build done.'
