package redis

import "github.com/redis/go-redis/v9"

type clusterClient struct {
	simpleClient
}

func OpenCluster(addresses []string, username, password string) Redis {
	return &clusterClient{simpleClient{c: redis.NewClusterClient(&redis.ClusterOptions{Addrs: addresses, Username: username, Password: password})}}
}
