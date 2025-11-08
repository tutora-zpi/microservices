package factory

import (
	"recorder-service/internal/domain/client"
)

type ClientFactory interface {
	CreateNewClient() client.Client
}
