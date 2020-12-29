package message

import (
	"context"
	"fmt"
)

// EmployeePayConfirmed send mesage to zeebe engine
func EmployeePayConfirmed(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("EmployeePayConfirmed").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		return err
	}
	return nil
}

// PackageLoaded send mesage to zeebe engine
func PackageLoaded(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("PackageLoaded").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// VehicleStarted send mesage to zeebe engine
func VehicleStarted(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("VehicleStarted").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// VehicleArrived send mesage to zeebe engine
func VehicleArrived(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("VehicleArrived").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// PackageUnloaded send mesage to zeebe engine
func PackageUnloaded(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("PackageUnloaded").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// ShipperReceived send mesage to zeebe engine
func ShipperReceived(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperReceived").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// ShipperCalled send mesage to zeebe engine
func ShipperCalled(orderID uint, timeConfirmed int64) error {
	variables := make(map[string]interface{})
	variables["time_confirmed "] = timeConfirmed
	ctx := context.Background()
	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperCalled").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
	_, err = request.Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// ShipperReceivedMoney send mesage to zeebe engine
func ShipperReceivedMoney(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperReceivedMoney").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// ShipperShipped send mesage to zeebe engine
func ShipperShipped(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperShipped").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// ShipperConfirmed send mesage to zeebe engine
func ShipperConfirmed(orderID uint, shipperConfirmed string) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperConfirmed").CorrelationKey(fmt.Sprint(orderID)).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}
