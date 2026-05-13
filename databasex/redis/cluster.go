package redis

import "github.com/redis/go-redis/v9"

type clusterClient struct {
	simpleClient
}

type ClusterOptions redis.ClusterOptions

func OpenCluster(addresses []string, withOptions ...func(options *ClusterOptions)) Redis {
	var options = &redis.ClusterOptions{Addrs: addresses}

	for i := range withOptions {
		withOptions[i]((*ClusterOptions)(options))
	}

	return &clusterClient{
		simpleClient{
			client: redis.NewClusterClient(options),
		},
	}
}
