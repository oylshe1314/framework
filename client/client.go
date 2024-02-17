package client

type Client interface {
	Init() error
	Close() error
}

type AsyncClient interface {
	Client
	Work() error
}

func InitClients(clients ...Client) (err error) {
	for _, client := range clients {
		err = client.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
