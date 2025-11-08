package factoryimpl

import (
	"os"
	"recorder-service/internal/app/interfaces/factory"
	"recorder-service/internal/config"
	"recorder-service/internal/domain/client"
	"recorder-service/internal/infrastructure/bus"
	wsclient "recorder-service/internal/infrastructure/ws_client"
)

type clientFactoryImpl struct {
	dispatcher bus.Dispachable
}

func (c *clientFactoryImpl) CreateNewClient() client.Client {
	url := os.Getenv(config.WS_GATEWAY_URL)
	return wsclient.NewWSClient(url, c.dispatcher)
}

func NewClientFactory(dispatcher bus.Dispachable) factory.ClientFactory {
	return &clientFactoryImpl{dispatcher: dispatcher}
}
