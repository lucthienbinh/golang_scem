package message

import (
	"context"
	"fmt"
)

// CustomerReceiveConfirmed send mesage to zeebe engine
func CustomerReceiveConfirmed(orderID uint) error {
	variables := make(map[string]interface{})
	variables["cus_receive_confirmed"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("CustomerReceiveConfirmed").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		return err
	}

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}
