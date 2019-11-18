// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package fake_test

import (
	"testing"

	"github.com/q191201771/naza/pkg/fake"

	"github.com/q191201771/naza/pkg/assert"
)

func TestNewWriter(t *testing.T) {
	_ = fake.NewWriter(fake.WriterTypeDoNothing)
}

func TestWriter_Write(t *testing.T) {
	var (
		w   *fake.Writer
		n   int
		err error
		b   = []byte("hello")
	)

	w = fake.NewWriter(fake.WriterTypeDoNothing)
	n, err = w.Write(b)
	assert.Equal(t, 5, n)
	assert.Equal(t, nil, err)

	w = fake.NewWriter(fake.WriterTypeReturnError)
	n, err = w.Write(b)
	assert.Equal(t, 0, n)
	assert.Equal(t, fake.ErrFakeWriter, err)

	w = fake.NewWriter(fake.WriterTypeIntoBuffer)
	n, err = w.Write(b)
	assert.Equal(t, 5, n)
	assert.Equal(t, nil, err)
}

func TestWriter_SetSpecificType(t *testing.T) {
	var (
		w   *fake.Writer
		n   int
		err error
		b   = []byte("hello")
	)
	w = fake.NewWriter(fake.WriterTypeDoNothing)
	w.SetSpecificType(map[uint32]fake.WriterType{
		0: fake.WriterTypeReturnError,
		2: fake.WriterTypeReturnError,
		4: fake.WriterTypeDoNothing,
	})

	expectedLen := map[int]int{
		0: 0,
		1: 5,
		2: 0,
		3: 5,
		4: 5,
		5: 5,
	}
	expectedErr := map[int]error{
		0: fake.ErrFakeWriter,
		1: nil,
		2: fake.ErrFakeWriter,
		3: nil,
		4: nil,
		5: nil,
	}

	for i := 0; i < 6; i++ {
		n, err = w.Write(b)
		assert.Equal(t, expectedLen[i], n)
		assert.Equal(t, expectedErr[i], err)
	}
}
