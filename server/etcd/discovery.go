package etcd

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientV3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"im/config"
	"sync"
	"time"
)

// Discovery 服务发现
type Discovery struct {
	client    *clientV3.Client //    etcd client
	serverMap sync.Map
}

func NewDiscovery() (*Discovery, error) {
	client, err := clientV3.New(clientV3.Config{
		Endpoints:   config.Configs.ETCD.Endpoints,
		DialTimeout: time.Duration(config.Configs.ETCD.Timeout) * time.Second,
	})
	if err != nil {
		zap.S().Error("etcd err:", err)
		return nil, err
	}
	return &Discovery{client: client}, nil
}

// WatchService 初始化服务列表和监视
func (d *Discovery) WatchService(prefix string) error {
	//根据前缀获取现有的key
	resp, err := d.client.Get(context.TODO(), prefix, clientV3.WithPrefix())
	if err != nil {
		return err
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i]; v != nil {
			d.serverMap.Store(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
		}
	}
	d.watcher(prefix)
	// 监听前缀
	return nil
}

func (d *Discovery) watcher(prefix string) {
	rch := d.client.Watch(context.TODO(), prefix, clientV3.WithPrefix())
	zap.S().Debugf("监听前缀: %s ..\n", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				zap.S().Infof("修改或新增, key:%s, value:%s\n", string(ev.Kv.Key), string(ev.Kv.Value))
				d.serverMap.Store(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				zap.S().Infof("删除, key:%s, value:%s\n", string(ev.Kv.Key), string(ev.Kv.Value))
				d.serverMap.Delete(string(ev.Kv.Key))
			}
		}
	}
}

func (d *Discovery) Close() error {
	return d.client.Close()
}

// GetServices 获取服务列表
func (d *Discovery) GetServices() []string {
	addrs := make([]string, 0)
	d.serverMap.Range(func(key, value interface{}) bool {
		addrs = append(addrs, value.(string))
		return true
	})
	return addrs
}
