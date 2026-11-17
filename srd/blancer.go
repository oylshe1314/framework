package srd

import "context"

type BalancerClient struct {
}

func (this *BalancerClient) Init(ctx context.Context) error {
	return nil
}

func (this *BalancerClient) Close() error {
	return nil
}
