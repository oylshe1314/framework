package redis

import (
	std "github.com/redis/go-redis/v9"
)

type cluster struct {
	simple
}

func OpenCluster(addrs []string, username, password string) Redis {
	return &cluster{simple{client: std.NewClusterClient(&std.ClusterOptions{Addrs: addrs, Username: username, Password: password})}}
}
