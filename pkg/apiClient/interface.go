package apiClient

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/watch"
	"io"
)

// TODO: 可能要改？
type Interface interface {
	Post(resourceURL string, context []byte) io.ReadCloser
	Get(resourceURL string, context []byte) io.ReadCloser
	GetObject(name string) (int, config.ApiObject, error)
	Put(resourceURL string, context []byte) io.ReadCloser
	PutObject(object config.ApiObject) (int, error)
	Delete(resourceURL string, context []byte) io.ReadCloser
	GetAll() (objectList config.ApiObjectList, err error)
	WatchAll() (watch.Interface, error)
	Watch(name string) (watch.Interface, error)
}
