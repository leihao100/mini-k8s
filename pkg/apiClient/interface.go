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
	Put(resourceURL string, context []byte) io.ReadCloser
	PutObject(name string, object config.ApiObject) (int, error)
	Delete(resourceURL string, context []byte) io.ReadCloser
	GetAll() (objectList config.ApiObjectList, err error)
	WatchAll() (watch.Interface, error)
	Watch(name string) (watch.Interface, error)
}
