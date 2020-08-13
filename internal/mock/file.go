package mock

import (
	"errors"
	td "github.com/huangjiahua/tempdesk"
	"io"
	"sync"
)

type fileInternal struct {
	rw   sync.RWMutex
	meta map[string]interface{}
	data []byte
}

type File struct {
	pos  int64
	file *fileInternal
}

func (f *File) Read(p []byte) (n int, err error) {
	n, err = f.ReadAt(p, f.pos)
	f.pos += int64(n)
	return
}

func (f *File) Write(p []byte) (n int, err error) {
	n, err = f.WriteAt(p, f.pos)
	f.pos += int64(n)
	return
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	var curr int64

	switch whence {
	case io.SeekStart:
		curr = offset
	case io.SeekCurrent:
		curr = f.pos + offset
	case io.SeekEnd:
		f.file.rw.RLock()
		curr = int64(len(f.file.data)) + offset
		f.file.rw.RUnlock()
	default:
		panic("not supported whence value")
	}

	if curr < 0 {
		panic("cannot seek to position before 0")
	}

	f.pos = curr
	return curr, nil
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("tempdesk.internal.mock.File.ReadAt: negative offset")
	}

	f.file.rw.RLock()
	defer f.file.rw.Unlock()

	if off >= int64(len(f.file.data)) {
		return 0, io.EOF
	}
	n = copy(p, f.file.data[off:])
	if n < len(p) {
		err = io.EOF
	}
	return
}

func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("tempdesk.internal.mock.File.WriteAt: negative offset")
	}

	f.file.rw.Lock()
	defer f.file.rw.Unlock()

	if off+int64(len(p)) >= int64(len(f.file.data)) {
		f.file.data = append(f.file.data, make([]byte, off+int64(len(p))-int64(len(f.file.data)))...)
	}

	n = copy(f.file.data[off:], p)
	return
}

func (f *File) Perm() td.FilePermission {
	panic("implement me")
}

func (f *File) Meta(key string) (value string, ok bool) {
	f.file.rw.RLock()
	defer f.file.rw.RUnlock()
	v, ok := f.file.meta[key]
	if !ok {
		return
	}
	value, ok = v.(string)
	return
}

func (f *File) WriteMeta(key string, value string) (err error) {
	return f.WriteFileMeta(key, value)
}

func (f *File) FileMeta(key string) (value interface{}, ok bool) {
	f.file.rw.RLock()
	defer f.file.rw.RUnlock()
	value, ok = f.file.meta[key]
	return
}

func (f *File) WriteFileMeta(key string, value interface{}) (err error) {
	f.file.rw.Lock()
	defer f.file.rw.Unlock()
	f.file.meta[key] = value
	return nil
}

func (f *File) Truncate(pos int64, data []byte) (err error) {
	panic("implement me")
}

type FileService struct {
}

func (fs FileService) File(path string) (err error) {
	panic("implement me")
}

func (fs FileService) Open(path string, flags int, perm td.FilePermission) (file td.File, err error) {
	panic("implement me")
}

func (fs FileService) Rename(dest string, src string) (err error) {
	panic("implement me")
}

func (fs FileService) Remove(path string) (err error) {
	panic("implement me")
}

func NewFileService() *FileService {
	return &FileService{}
}
