// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package circularqueue_test

import (
	"testing"

	"github.com/q191201771/naza/pkg/assert"
	"github.com/q191201771/naza/pkg/circularqueue"
)

func TestCircularQueue(t *testing.T) {
	var (
		err error
		n   int
		b   bool
		v   interface{}
	)

	q := circularqueue.New(3)
	assert.IsNotNil(t, q)

	// empty
	_, err = q.PopFront()
	assert.IsNotNil(t, err)
	_, err = q.Front()
	assert.IsNotNil(t, err)
	_, err = q.Back()
	assert.IsNotNil(t, err)
	_, err = q.At(0)
	assert.IsNotNil(t, err)
	n = q.Size()
	assert.Equal(t, 0, n)
	b = q.Full()
	assert.Equal(t, false, b)
	b = q.Empty()
	assert.Equal(t, true, b)

	err = q.PushBack(1)
	assert.Equal(t, nil, err)

	// [1]
	v, err = q.Front()
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, v.(int))
	v, err = q.Back()
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, v.(int))
	v, err = q.At(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, v.(int))
	n = q.Size()
	assert.Equal(t, 1, n)
	b = q.Full()
	assert.Equal(t, false, b)
	b = q.Empty()
	assert.Equal(t, false, b)

	err = q.PushBack(2)
	assert.Equal(t, nil, err)
	err = q.PushBack(3)
	assert.Equal(t, nil, err)

	// [1, 2, 3]
	v, err = q.Front()
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, v.(int))
	v, err = q.Back()
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, v.(int))
	v, err = q.At(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, v.(int))
	v, err = q.At(1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, v.(int))
	v, err = q.At(2)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, v.(int))
	n = q.Size()
	assert.Equal(t, 3, n)
	b = q.Full()
	assert.Equal(t, true, b)
	b = q.Empty()
	assert.Equal(t, false, b)

	v, err = q.PopFront()
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, v.(int))

	err = q.PushBack(400)
	assert.Equal(t, nil, err)

	err = q.PushBack(500)
	assert.IsNotNil(t, err)

	// [2, 3, 400]
	v, err = q.Front()
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, v.(int))
	v, err = q.Back()
	assert.Equal(t, nil, err)
	assert.Equal(t, 400, v.(int))
	v, err = q.At(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, v.(int))
	v, err = q.At(1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, v.(int))
	v, err = q.At(2)
	assert.Equal(t, nil, err)
	assert.Equal(t, 400, v.(int))
	n = q.Size()
	assert.Equal(t, 3, n)
	b = q.Full()
	assert.Equal(t, true, b)
	b = q.Empty()
	assert.Equal(t, false, b)
}
