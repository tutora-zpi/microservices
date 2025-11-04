package factoryimpl

import (
	"os"
	"recorder-service/internal/app/interfaces/factory"
	"recorder-service/internal/config"
	"recorder-service/internal/domain/client"
	ws "recorder-service/internal/infrastructure/ws_client"
)

type clientFactoryImpl struct {
}

// CreateNewClient implements factory.ClientFactory.
func (c *clientFactoryImpl) CreateNewClient() client.Client {
	url := os.Getenv(config.WS_GATEWAY_URL)
	return ws.NewSocketClient(url, nil)
}

func NewClientFactory() factory.ClientFactory {
	return &clientFactoryImpl{}
}
