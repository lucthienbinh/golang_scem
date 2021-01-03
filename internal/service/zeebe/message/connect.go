package message

import (
	"os"

	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

var (
	zbClient zbc.Client
)

// ConnectZeebeEngine function
func ConnectZeebeEngine() error {
	gatewayAddress := os.Getenv("BROKER_ADDRESS")
	newZbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddress,
		UsePlaintextConnection: true,
	})

	if err != nil {
		return err
	}

	zbClient = newZbClient
	return nil
}
