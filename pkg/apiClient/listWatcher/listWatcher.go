package listWatcher

import (
	"MiniK8S/pkg/apiClient"
	"io"
)

type ListWatcher struct {
	client apiClient.Client
}

func New(client apiClient.Client) *ListWatcher {
	return &ListWatcher{client: client}
}

func (lw *ListWatcher) Watch() io.ReadCloser {
	url := lw.client.BuildURL("watch")
	res := lw.client.Get(url, nil)
	return res
}

func (lw *ListWatcher) List() io.ReadCloser {
	url := lw.client.BuildURL("get")
	res := lw.client.Get(url, nil)
	return res
}
