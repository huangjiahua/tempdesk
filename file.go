package tempdesk

type FileMeta map[string]string

type FileData interface {
}

type File struct {
	Path string
	Name string
	Meta FileMeta
	Data FileData
}

type FileService interface {
}
