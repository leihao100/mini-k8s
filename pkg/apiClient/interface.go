package apiClient

import (
	"MiniK8S/pkg/api/config"
	"MiniK8S/pkg/api/watch"
	"io"
)

type Interface interface {
	URL() string
	BuildURL(requestType RequestType) string
	Post(resourceURL string, context []byte) io.ReadCloser
	Get(resourceURL string, context []byte) io.ReadCloser
	Put(resourceURL string, context []byte) io.ReadCloser
	Delete(resourceURL string, context []byte) io.ReadCloser
	GetAll() (objectList config.ApiObjectList, err error)
	WatchAll() (watch.Interface, error)
	Watch(name string) (watch.Interface, error)
}
