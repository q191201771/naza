// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/q191201771/naza/pkg/filebatch"
	"github.com/q191201771/naza/pkg/nazalog"
)

// 帮助找出源码中多个大写字母连接在一起的地方

func main() {
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		//option.LevelFlag = false
		option.TimestampWithMsFlag = false
		option.TimestampFlag = false
		//option.ShortFileFlag = false
	})

	dir := parseFlag()

	// 遍历所有go文件
	err := filebatch.Walk(dir, true, ".go", func(path string, info os.FileInfo, content []byte, err error) []byte {
		if err != nil {
			nazalog.Warnf("read file failed. file=%s, err=%+v", path, err)
			return nil
		}

		nazalog.Tracef("path:%s", path)

		//checkModFile(path, content)
		//return nil

		//免检的文件：
		if strings.Contains(path, "/pkg/alpha/stun/") {
			return nil
		}

		// 免检文件：测试文件
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		lines := bytes.Split(content, []byte{'\n'})

		// 免检的行：
		ignContainsKeyList := []string{
			// 字符串
			"\"",
			// 16进制的数字
			"0x",
			// 接口
			"IBufWriter",
			"IClientSession",
			"IServerSession",
			"IClientSessionLifecycle",
			"IServerSessionLifecycle",
			"ISessionStat",
			"ISessionUrlContext",
			"IObject",
			"IPathStrategy",
			"IPathRequestStrategy",
			"IPathWriteStrategy",
			"IQueueObserver",
			"IHandshakeClient",
			"IRtpUnpacker",
			"IRtpUnpackContainer",
			"IRtpUnpackerProtocol",
			"IInterleavedPacketWriter",
			"filesystemlayer.IFileSystemLayer",
			"filesystemlayer.IFile",
			"IFile",
			"IFileSystemLayer",
			// RTSP相关
			"HeaderCSeq",
			"ARtpMap",
			"AFmtPBase",
			"AControl",
			// 标准库
			".URL",
			".TLS",
			".SIGUSR",
			"ServeHTTP(",
			".URI",
			".RequestURI",
			"io.EOF",
			"net.UDPAddr",
			"net.UDPConn",
			"net.ResolveUDPAddr",
			"net.ListenUDP",
			"WriteToUDP",
			"ReadFromUDP",
			"time.RFC1123",
			"runtime.GOOS",
			"crc32.ChecksumIEEE",
			"cipher.NewCBCEncrypter",
			"cipher.NewCBCDecrypter",
			"os.O_CREATE",
			//
			"LAddr",
			"RAddr",
		}

		// 注释
		ignPrefixKeyList := []string{
			"//",
			"/*",
		}

		// 逐行分析
		for j, line := range lines {
			// 免检的行：
			if j == 3 && strings.Contains(string(line), "MIT-style license") {
				continue
			}

			ignFlag := false
			for _, k := range ignContainsKeyList {
				if strings.Contains(string(line), k) {
					nazalog.Debugf("ign contains line:%s %s", string(line), k)
					ignFlag = true
				}
			}
			if ignFlag {
				continue
			}

			for _, k := range ignPrefixKeyList {
				if strings.HasPrefix(strings.TrimSpace(string(line)), k) {
					nazalog.Debugf("ign prefix line:%s %s", string(line), k)
					ignFlag = true
				}
			}
			if ignFlag {
				continue
			}

			for i := range line {
				// 连续两个字符是大写字母
				if i == 0 {
					continue
				}
				if isCap(line[i]) && isCap(line[i-1]) {
					nazalog.Infof("%s:%d %s", path, j+1, string(highlightSerialCap(line)))
					break
				}
			}
		}
		return nil
	})
	nazalog.Assert(nil, err)
}

// 是否大写字母
func isCap(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 有连续大写的地方高亮显示
func highlightSerialCap(line []byte) []byte {
	var ret []byte
	var cache []byte
	for i := range line {
		if isCap(line[i]) {
			cache = append(cache, line[i])
		} else {
			if cache != nil {
				if len(cache) > 1 {
					ret = append(ret, []byte("\033[22;31m")...)
				}
				ret = append(ret, cache...)
				if len(cache) > 1 {
					ret = append(ret, []byte("\033[0m")...)
				}
				cache = nil
			}
			ret = append(ret, line[i])
		}
	}
	if cache != nil {
		if len(cache) > 1 {
			ret = append(ret, []byte("\033[22;31m")...)
		}
		ret = append(ret, cache...)
		if len(cache) > 1 {
			ret = append(ret, []byte("\033[0m")...)
		}
	}

	return ret
}

// 一段检查文件修改后和修改前的逻辑，修改是否符合预期
func checkModFile(path string, content []byte) {
	beforeModPath := "x"
	afterModPath := "x"

	ignFileList := []string{
		"pkg/rtmp/server_session.go",
		"pkg/base/websocket.go",
		"pkg/rtsp/client_command_session.go",
	}
	for _, f := range ignFileList {
		if strings.HasSuffix(path, f) {
			return
		}
	}
	beforeModFilename := strings.ReplaceAll(path, afterModPath, beforeModPath)
	beforeModContent, err := ioutil.ReadFile(beforeModFilename)
	nazalog.Assert(nil, err)

	// 理论上，大部分修改，不影响文件大小
	if len(content) != len(beforeModContent) {
		nazalog.Errorf("file size not match. path=%s, len(b=%d, a=%d)", path, len(content), len(beforeModContent))
	}

	notEqualFlag := false
	// 不管大小是否相等，取最小值，逐个字节比较内容
	// 理论上，大部分修改，要么是相等，要么是将新内容从小写转换回大写就相等
	for i := 0; i < min(len(content), len(beforeModContent)); i++ {
		if content[i] != beforeModContent[i] && content[i]-32 != beforeModContent[i] {
			nazalog.Errorf("-----a-----\n%s\n-----b-----\n%s", string(content[i:i+128]), string(beforeModContent[i:i+128]))
			notEqualFlag = true
			break
		}
	}
	if notEqualFlag {
		nazalog.Errorf("%s", path)
	}
}

func parseFlag() string {
	dir := flag.String("d", "", "dir of source")
	flag.Parse()
	if *dir == "" {
		flag.Usage()
		os.Exit(1)
	}
	return *dir
}
