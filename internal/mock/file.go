package mock

import (
	td "github.com/huangjiahua/tempdesk"
)

type FileService struct {
}

func (fs FileService) File(path string) (err error) {
	panic("implement me")
}

func (fs FileService) Open(path string, flags int, perm *td.FilePermission) (file *td.File, err error) {
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
