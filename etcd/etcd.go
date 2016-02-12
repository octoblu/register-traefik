package etcd

import (
	"time"

	"github.com/octoblu/go-simple-etcd-client/etcdclient"
)

// Del deletes a key from etcd
func Del(uri, key string) error {
	client, err := etcdclient.Dial(uri)
	if err != nil {
		return err
	}

	return client.Del(key)
}

// DelDir deletes directory from etcd
func DelDir(uri, key string) error {
	client, err := etcdclient.Dial(uri)
	if err != nil {
		return err
	}

	return client.DelDir(key)
}

// Set sets a key on etcd
func Set(uri, key, value string) error {
	client, err := etcdclient.Dial(uri)
	if err != nil {
		return err
	}

	return client.Set(key, value)
}

// UpdateDirWithTTL updates the ttl on the dir
func UpdateDirWithTTL(uri, key string, ttlSeconds int) error {
	client, err := etcdclient.Dial(uri)
	if err != nil {
		return err
	}

	ttl := time.Duration(ttlSeconds) * time.Second
	return client.UpdateDirWithTTL(key, ttl)
}
