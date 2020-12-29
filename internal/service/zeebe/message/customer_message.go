package message

import (
	"context"
	"fmt"
)

// CustomerPayConfirmed send mesage to zeebe engine
func CustomerPayConfirmed(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("CustomerPayConfirmed").CorrelationKey(fmt.Sprintf("%x", orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// CustomerReceiveConfirmed send mesage to zeebe engine
func CustomerReceiveConfirmed(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("CustomerReceiveConfirmed").CorrelationKey(fmt.Sprintf("%x", orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}
