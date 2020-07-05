package storage

type PutterGetter interface {
	Put(name string, value interface{}) (err error)
	Get(name string) (value interface{}, err error)
}
