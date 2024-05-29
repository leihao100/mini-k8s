package etcd

import (
	"MiniK8S/config"
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var etcdEndpoint = "http://" + config.EtcdHost() + config.EtcdPort()
var etcdConfig clientv3.Config
var etcdClient *clientv3.Client
var requestTimeout = time.Second
var err error
var RVM ResourceVersionManager
var VersionLock sync.RWMutex

type ResourceVersionManager struct {
	version int64
	mutex   sync.RWMutex
}

func (resourceVersionManager *ResourceVersionManager) Init(version int64) {
	resourceVersionManager.mutex.Lock()
	resourceVersionManager.version = version
	resourceVersionManager.mutex.Unlock()
}

func (resourceVersionManager *ResourceVersionManager) SetResourceVersion(version int64) {
	resourceVersionManager.mutex.Lock()
	if version > resourceVersionManager.version {
		resourceVersionManager.version = version
	}
	resourceVersionManager.mutex.Unlock()
}

func (resourceVersionManager *ResourceVersionManager) GetResourceVersion() int64 {
	resourceVersionManager.mutex.RLock()
	defer resourceVersionManager.mutex.RUnlock()
	return resourceVersionManager.version
}

func (resourceVersionManager *ResourceVersionManager) GetNextResourceVersion() int64 {
	resourceVersionManager.mutex.RLock()
	defer resourceVersionManager.mutex.RUnlock()
	return resourceVersionManager.version + 1
}

func Init() {
	etcdConfig = clientv3.Config{
		Endpoints:            []string{etcdEndpoint},
		DialTimeout:          60 * time.Second,
		DialKeepAliveTimeout: 60 * time.Second,
	}
	etcdClient, err = clientv3.New(etcdConfig)
	if err != nil {
		fmt.Printf("[etcd]%v", err)
	} else {
		fmt.Printf("[etcd] connect to etcd success\n")
	}
	//for test
	fmt.Println("[debug] etcd endpoint:", etcdEndpoint)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*10)
	status, err := etcdClient.Status(ctx, etcdEndpoint)
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
	}
	RVM.Init(status.Header.Revision)
}

func Close() {
	err = etcdClient.Close()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
	} else {
		fmt.Printf("[etcd] etcd client has closed\n")
	}
}

func Put(key string, value string) (version int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Put(ctx, key, value)
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
	}
	version = res.Header.Revision
	RVM.SetResourceVersion(version)
	return version, err
}

func Get(key string) (value string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Get(ctx, key)
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
		return "", err
	}
	if res.Count > 0 {
		return string(res.Kvs[0].Value), err
	} else {
		return "", err
	}
}

func GetAllWithPrefix(key string) (values []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Get(ctx, key, clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
		return nil, err
	}
	for _, kv := range res.Kvs {
		values = append(values, string(kv.Value))
	}
	return values, err
}

func Exist(key string) (value bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Get(ctx, key)
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
		return false, err
	}
	if res.Count == 0 {
		return false, err
	} else {
		return true, err
	}
}

func Delete(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Delete(ctx, key)
	defer cancel()
	RVM.SetResourceVersion(res.Header.Revision)
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
	}
	return err
}

func DeleteAllWithPrefix(key string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Delete(ctx, key, clientv3.WithPrefix())
	defer cancel()
	RVM.SetResourceVersion(res.Header.Revision)

	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
	}
	return err
}

func Clear() (err error) {
	return DeleteAllWithPrefix("")
}

func PutWithVersion(key string, value string, oldVersion int64) (newVersion int64, success bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Put(ctx, key, value, clientv3.WithPrevKV())
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
	}
	newVersion = res.Header.Revision
	RVM.SetResourceVersion(newVersion)
	if oldVersion != res.PrevKv.ModRevision {
		fmt.Printf("[etcd] OldVersion %v and res.PrevKv.ModRevision %v mismatch\n", oldVersion, res.PrevKv.ModRevision)
		return newVersion, false, err
	}
	return newVersion, true, err
}

func GetWithVersion(key string) (value string, version int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Get(ctx, key)
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
		return "", version, err
	}
	if res.Count > 0 {
		version = res.Kvs[0].ModRevision
		return string(res.Kvs[0].Value), version, err
	} else {
		return "", version, err
	}
}

func ExistWithVersion(key string) (value bool, version int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	res, err := etcdClient.Get(ctx, key)
	defer cancel()
	if err != nil {
		fmt.Printf("[etcd]%v\n", err)
		return false, version, err
	}
	if res.Count == 0 {
		return false, version, err
	} else {
		version = res.Kvs[0].ModRevision
		return true, version, err
	}
}

func Watch(key string) (context.CancelFunc, chan *clientv3.Event) {
	ctx, cancel := context.WithCancel(context.Background())
	rch := etcdClient.Watch(ctx, key, clientv3.WithPrevKV())
	ch := make(chan *clientv3.Event)
	go func(rch clientv3.WatchChan, ch chan *clientv3.Event) {
		for res := range rch {
			for _, ev := range res.Events {
				ch <- (*clientv3.Event)(ev)
			}
		}
	}(rch, ch)
	return cancel, ch
}

func WatchAllWithPrefix(key string) (context.CancelFunc, chan *clientv3.Event) {
	ctx, cancel := context.WithCancel(context.Background())
	rch := etcdClient.Watch(ctx, key, clientv3.WithPrefix(), clientv3.WithPrevKV())
	ch := make(chan *clientv3.Event)
	go func(rch clientv3.WatchChan, ch chan *clientv3.Event) {
		for res := range rch {
			for _, ev := range res.Events {
				ch <- (*clientv3.Event)(ev)
			}
		}
	}(rch, ch)
	return cancel, ch
}
