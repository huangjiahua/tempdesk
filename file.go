package tempdesk

import "io"

const (
	FileNotExist = "file not exists"
)

type FilePermission interface {
	AllowUser(name string)
	BlockUser(name string)
	AllowUserMeta(key, value string)
	BlockUserMeta(key, value string)
	AllowAllUser()
	BlockAllUser()
	AllowPublic(code string)
	AllowCode(code string)
	BlockPublic()
	BlockCode(code string)

	TestUser(user User) bool
	TestCode(code string) bool
}

type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.ReaderAt
	io.WriterAt

	Perm() FilePermission

	Meta(key string) (value string, ok bool)
	WriteMeta(key string, value string) (err error)
	FileMeta(key string) (value interface{}, ok bool)
	WriteFileMeta(key string, value interface{}) (err error)

	Truncate(pos int64, data []byte) (err error)
}

type FileService interface {
	File(path string) (err error)
	Open(path string, flags int, perm FilePermission) (file File, err error)
	Rename(dest string, src string) (err error)
	Remove(path string) (err error)
}

type FileServiceError struct {
	Kind string
	Err  error
}

func (f *FileServiceError) Error() string {
	return f.Kind
}
