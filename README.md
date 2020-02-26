<p align="center">
<br>
Go语言基础库
<br><br>
<a title="TravisCI" target="_blank" href="https://www.travis-ci.org/q191201771/naza"><img src="https://www.travis-ci.org/q191201771/naza.svg?branch=master"></a>
<a title="codecov" target="_blank" href="https://codecov.io/gh/q191201771/naza"><img src="https://codecov.io/gh/q191201771/naza/branch/master/graph/badge.svg?style=flat-square"></a>
<a title="goreportcard" target="_blank" href="https://goreportcard.com/report/github.com/q191201771/naza"><img src="https://goreportcard.com/badge/github.com/q191201771/naza?style=flat-square"></a>
<br>
<a title="codeline" target="_blank" href="https://github.com/q191201771/naza"><img src="https://sloc.xyz/github/q191201771/naza/?category=code"></a>
<a title="license" target="_blank" href="https://github.com/q191201771/naza/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
<a title="lastcommit" target="_blank" href="https://github.com/q191201771/naza/commits/master"><img src="https://img.shields.io/github/commit-activity/m/q191201771/naza.svg?style=flat-square"></a>
<a title="commitactivity" target="_blank" href="https://github.com/q191201771/naza/graphs/commit-activity"><img src="https://img.shields.io/github/last-commit/q191201771/naza.svg?style=flat-square"></a>
<br>
<a title="pr" target="_blank" href="https://github.com/q191201771/naza/pulls"><img src="https://img.shields.io/github/issues-pr-closed/q191201771/naza.svg?style=flat-square&color=FF9966"></a>
<a title="hits" target="_blank" href="https://github.com/q191201771/naza"><img src="https://hits.b3log.org/q191201771/naza.svg?style=flat-square"></a>
<a title="language" target="_blank" href="https://github.com/q191201771/naza"><img src="https://img.shields.io/github/languages/count/q191201771/naza.svg?style=flat-square"></a>
<a title="toplanguage" target="_blank" href="https://github.com/q191201771/naza"><img src="https://img.shields.io/github/languages/top/q191201771/naza.svg?style=flat-square"></a>
<a title="godoc" target="_blank" href="https://godoc.org/github.com/q191201771/naza"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square"></a>
<br><br>
<a title="watcher" target="_blank" href="https://github.com/q191201771/naza/watchers"><img src="https://img.shields.io/github/watchers/q191201771/naza.svg?label=Watchers&style=social"></a>&nbsp;&nbsp;
<a title="star" target="_blank" href="https://github.com/q191201771/naza/stargazers"><img src="https://img.shields.io/github/stars/q191201771/naza.svg?label=Stars&style=social"></a>&nbsp;&nbsp;
<a title="fork" target="_blank" href="https://github.com/q191201771/naza/network/members"><img src="https://img.shields.io/github/forks/q191201771/naza.svg?label=Forks&style=social"></a>&nbsp;&nbsp;
</p>

---

#### 工程目录说明

```
pkg/                    ...... 源码包
    |-- bininfo/        ...... 将编译时源码的 git 版本信息（当前 commit log 的 sha 值和 commit message），编译时间，Go 版本，平台打入程序中
    |-- snowflake/      ...... 分布式唯一性 ID 生成器
    |-- bitrate/        ...... 计算带宽
    |-- consistenthash/ ...... 一致性哈希
    |-- nazalog/        ...... 日志库
    |-- slicebytepool/  ...... []byte 内存池
    |-- assert/         ...... 提供了单元测试时的断言功能，减少一些模板代码
    |-- nazastring/     ...... string 和 []byte 相关的操作
    |-- bele/           ...... 提供了大小端的转换操作
    |-- nazaatomic/     ...... 原子操作
    |-- nazajson/       ...... json 操作
    |-- taskpool/       ...... 协程池
    |-- connection/     ...... 对 net.Conn 接口的二次封装
    |-- filebatch/      ...... 文件批处理操作
    |-- nazamd5/        ...... md5 操作
    |-- ratelimit/      ...... 限流器，令牌桶，漏桶
    |-- ic/             ...... 将整型切片压缩成二进制字节切片
    |-- unique/         ...... 对象唯一 ID
    |-- fake/           ...... stub和mock相关，实现一些常用的接口，辅助测试其它代码
playground/             ...... Go实验代码片段
demo/                   ...... 示例相关的代码
bin/                    ...... 可执行文件编译输出目录
```

#### 依赖

无任何第三方依赖

#### 项目名 naza 由来

本仓库主要用于存放我自己写的一些 Go 基础库代码。目前主要服务于我的另一个项目： [lal](https:////github.com/q191201771/lal)

naza 即哪吒（正确拼音为 nezha，我女儿发音读作 naza，少一个字母，挺好~），希望本仓库以后能像三头六臂，有多种武器的哪吒一样，为我提供一个趁手的工具箱。

#### 相关文档

* pkg/snowflake [分布式ID生成算法snowflake介绍及Go语言实现](https://pengrl.com/p/20041/)
* pkg/bininfo [给Go程序加入编译版本时间等信息](https://pengrl.com/p/37397/)
* pkg/nazastring [Go语言中[]byte和string类型相互转换时的性能分析和优化](https://www.pengrl.com/p/31544/)

* [Go创建对象时，如何优雅的传递初始化参数](https://pengrl.com/p/60015/)
* playground/p3 [利用CPU cache特性优化Go程序](https://pengrl.com/p/9125/)
* playground/p4 [老弟在吗，我的Go标准库中的二分查找有bug！](https://pengrl.com/p/20011/)

#### 联系我

欢迎扫码加我微信，进行技术交流或扯淡。

<img src="https://pengrl.com/images/yoko_vx.jpeg" width="180" height="180" />
