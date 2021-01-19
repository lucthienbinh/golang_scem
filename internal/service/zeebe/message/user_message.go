package message

import (
	"context"
	"fmt"
	"time"
)

// PaymentConfirmed send mesage to zeebe engine
func PaymentConfirmed(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("PaymentConfirmed").CorrelationKey(fmt.Sprint(orderID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		return err
	}
	return nil
}

// PackageLoaded send mesage to zeebe engine
func PackageLoaded(longShipID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("PackageLoaded").CorrelationKey(fmt.Sprint(longShipID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// VehicleStarted send mesage to zeebe engine
func VehicleStarted(longShipID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("VehicleStarted").CorrelationKey(fmt.Sprint(longShipID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// VehicleArrived send mesage to zeebe engine
func VehicleArrived(longShipID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("VehicleArrived").CorrelationKey(fmt.Sprint(longShipID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// PackageUnloaded send mesage to zeebe engine
func PackageUnloaded(longShipID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("PackageUnloaded").CorrelationKey(fmt.Sprint(longShipID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// ShipperCalled send mesage to zeebe engine
func ShipperCalled(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperCalled").CorrelationKey(fmt.Sprint(orderID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		return nil
	}
	return nil
}

// ShipperReceivedMoney send mesage to zeebe engine
func ShipperReceivedMoney(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperReceivedMoney").CorrelationKey(fmt.Sprint(orderID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// ShipperShipped send mesage to zeebe engine
func ShipperShipped(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperShipped").CorrelationKey(fmt.Sprint(orderID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}

// ShipperConfirmed send mesage to zeebe engine
func ShipperConfirmed(orderID uint) error {
	ctx := context.Background()
	_, err := zbClient.NewPublishMessageCommand().MessageName("ShipperConfirmed").CorrelationKey(fmt.Sprint(orderID)).TimeToLive(1 * time.Minute).Send(ctx)
	if err != nil {
		// failed to set the updated variables
		return err
	}
	return nil
}
