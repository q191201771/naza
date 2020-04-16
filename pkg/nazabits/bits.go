// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package nazabits

// @param pos: 取值范围 [0, 7]，0表示最低位
func GetBit8(v uint8, pos int) int {
	return GetBits8(v, pos, 1)
}

// @param pos: 取值范围 [0, 7]，0表示最低位
// @param n:   取多少位， 取值范围 [1, 8]
//
// 举例，GetBits8(105, 2, 4) = 10（即1010）
//   v: 0110 1001
// pos:       2
//   n:   .. ..
//
func GetBits8(v uint8, pos int, n int) int {
	m := []uint8{0, 1, 3, 7, 15, 31, 63, 127, 255}
	return int(v >> uint(pos) & m[n])
}
