package httpx

import "context"

type HttpClient struct {
}

func (this *HttpClient) Init(ctx context.Context) error {
	return nil
}

func (this *HttpClient) Close() error {
	return nil
}
