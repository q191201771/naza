// Copyright 2021, Chef.  All rights reserved.
// https://github.com/q191201771/naza
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package filesystemlayer

import (
	"io/ioutil"
	"os"
)

type FSLDisk struct {
}

func (f *FSLDisk) Type() FSLType {
	return FSLTypeDisk
}

func (f *FSLDisk) Create(name string) (IFile, error) {
	return os.Create(name)
}

func (f *FSLDisk) Rename(oldpath string, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (f *FSLDisk) MkdirAll(path string, perm uint32) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

func (f *FSLDisk) Remove(name string) error {
	return os.Remove(name)
}

func (f *FSLDisk) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f *FSLDisk) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (f *FSLDisk) WriteFile(filename string, data []byte, perm uint32) error {
	return ioutil.WriteFile(filename, data, os.FileMode(perm))
}
