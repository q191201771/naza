<p align="center">
<br>
Go语言基础库
<br><br>
<a title="TravisCI" target="_blank" href="https://www.travis-ci.org/q191201771/nezha"><img src="https://www.travis-ci.org/q191201771/nezha.svg?branch=master"></a>
<a title="codecov" target="_blank" href="https://codecov.io/gh/q191201771/nezha"><img src="https://codecov.io/gh/q191201771/nezha/branch/master/graph/badge.svg?style=flat-square"></a>
<a title="goreportcard" target="_blank" href="https://goreportcard.com/report/github.com/q191201771/nezha"><img src="https://goreportcard.com/badge/github.com/q191201771/nezha?style=flat-square"></a>
<br>
<a title="codesize" target="_blank" href="https://github.com/q191201771/nezha"><img src="https://img.shields.io/github/languages/code-size/q191201771/nezha.svg?style=flat-square?style=flat-square"></a>
<a title="license" target="_blank" href="https://github.com/q191201771/nezha/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
<a title="lastcommit" target="_blank" href="https://github.com/q191201771/nezha/commits/master"><img src="https://img.shields.io/github/commit-activity/m/q191201771/nezha.svg?style=flat-square"></a>
<a title="commitactivity" target="_blank" href="https://github.com/q191201771/nezha/graphs/commit-activity"><img src="https://img.shields.io/github/last-commit/q191201771/nezha.svg?style=flat-square"></a>
<br>
<a title="pr" target="_blank" href="https://github.com/q191201771/nezha/pulls"><img src="https://img.shields.io/github/issues-pr-closed/q191201771/nezha.svg?style=flat-square&color=FF9966"></a>
<a title="hits" target="_blank" href="https://github.com/q191201771/nezha"><img src="https://hits.b3log.org/q191201771/nezha.svg?style=flat-square"></a>
<a title="language" target="_blank" href="https://github.com/q191201771/nezha"><img src="https://img.shields.io/github/languages/count/q191201771/nezha.svg?style=flat-square"></a>
<a title="toplanguage" target="_blank" href="https://github.com/q191201771/nezha"><img src="https://img.shields.io/github/languages/top/q191201771/nezha.svg?style=flat-square"></a>
<a title="godoc" target="_blank" href="https://godoc.org/github.com/q191201771/nezha"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square"></a>
<br><br>
<a title="watcher" target="_blank" href="https://github.com/q191201771/nezha/watchers"><img src="https://img.shields.io/github/watchers/q191201771/nezha.svg?label=Watchers&style=social"></a>&nbsp;&nbsp;
<a title="star" target="_blank" href="https://github.com/q191201771/nezha/stargazers"><img src="https://img.shields.io/github/stars/q191201771/nezha.svg?label=Stars&style=social"></a>&nbsp;&nbsp;
<a title="fork" target="_blank" href="https://github.com/q191201771/nezha/network/members"><img src="https://img.shields.io/github/forks/q191201771/nezha.svg?label=Forks&style=social"></a>&nbsp;&nbsp;
</p>

---

#### 工程目录说明

```
pkg/                  ......源码包
    |-- assert/       ......提供了单元测试时的断言功能，减少一些模板代码
    |-- bele/         ......提供了大小端的转换操作
    |-- bininfo/      ......将编译时的git版本号，时间，Go编译器信息打入程序中
    |-- connection/   ......对 net.Conn 接口的二次封装
    |-- log/          ......日志库
    |-- mockserver    ......模拟一些服务端，用于快速测试其它代码
    |-- mockwriter    ......模拟Writer接口，用于快速测试其它代码
    |-- unique/       ......对象唯一ID
demo/                 ......示例相关的代码
    |-- connstat/     ......简单测试 net.Conn.SetWriteDeadline 的性能
bin/                  ......可执行文件编译输出目录
```

#### 依赖

无任何第三方依赖

#### 项目名 nezha 由来

本仓库主要用于存放我自己写的一些 Go 基础库代码。目前只服务于我的另一个项目： [lal](https:////github.com/q191201771/lal)

nezha 即 哪吒，希望本仓库以后能像三头六臂，有多种武器的哪吒一样，为我提供多种工具。
