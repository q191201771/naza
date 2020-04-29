// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package p12_test

import (
	"log"
	"os"
	"runtime"
	"testing"
)

//BenchmarkRuntimeCaller-4    	 2417739	       488 ns/op	     216 B/op	       2 allocs/op
func BenchmarkRuntimeCaller(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runtime.Caller(0)
	}
}

//BenchmarkRuntimeCaller2-4   	 1213971	       983 ns/op	     216 B/op	       2 allocs/op
func BenchmarkRuntimeCaller2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runtime.Caller(2)
	}
}

//BenchmarkLog-4              	  754929	      1672 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLog(b *testing.B) {
	fp, _ := os.Create("/dev/null")
	log.SetOutput(fp)
	log.SetFlags(0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Printf("a")
	}
}

//BenchmarkLogWith-4          	  344067	      3403 ns/op	     216 B/op	       2 allocs/op
func BenchmarkLogWith(b *testing.B) {
	fp, _ := os.Create("/dev/null")
	log.SetOutput(fp)
	log.SetFlags(log.Lshortfile)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		log.Printf("a")
	}
}
