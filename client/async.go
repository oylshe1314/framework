package client

type AsyncClient interface {
	Client

	Work() error
}
