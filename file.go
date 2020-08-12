package tempdesk

import "io"

type FilePermission struct {
	owner       string
	exclusive   bool
	notAllowed  map[string]bool
	allowed     map[string]bool
	allowedMeta map[string]string
}

type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.ReaderAt
	io.WriterAt

	Perm() *FilePermission

	Meta(key string) (value string)
	WriteMeta(key string, value string) (err error)
	FileMeta(key string) (value interface{})
	WriteFileMeta(key string, value interface{}) (err error)

	Truncate(pos int64, data []byte) (err error)
}

type FileService interface {
	File(path string) (err error)
	Open(path string, flags int, perm *FilePermission) (file *File, err error)
	Rename(dest string, src string) (err error)
	Remove(path string) (err error)
}
