package zk

import "context"

type SubscriberClient struct {
	client
}

func (this *SubscriberClient) Init(ctx context.Context) error {
	var err = this.client.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}
