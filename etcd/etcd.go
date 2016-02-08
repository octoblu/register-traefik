package etcd

import "github.com/octoblu/go-simple-etcd-client/etcdclient"

// Del deletes an key from etcd
func Del(uri, key string) error {
	client, err := etcdclient.Dial(uri)
	if err != nil {
		return err
	}

	return client.Del(key)
}

// Set sets a key on etcd
func Set(uri, key, value string) error {
	client, err := etcdclient.Dial(uri)
	if err != nil {
		return err
	}

	return client.Set(key, value)
}
