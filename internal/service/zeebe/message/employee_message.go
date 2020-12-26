package message

import (
	"context"
	"fmt"
)

// MoneyReceived send mesage to zeebe engine
func MoneyReceived(payEmployeeID, orderID uint) error {
	variables := make(map[string]interface{})
	variables["pay_employee_id"] = payEmployeeID
	variables["pay_status"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("MoneyReceived").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		return err
	}

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		return err
	}
	return nil
}

// PaymentInfoUpdated send mesage to zeebe engine
func PaymentInfoUpdated(orderID uint) error {
	variables := make(map[string]interface{})

	request, err := zbClient.NewPublishMessageCommand().MessageName("PaymentInfoUpdated").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// PackageLoaded send mesage to zeebe engine
func PackageLoaded(orderID, emplLoadID uint) error {
	variables := make(map[string]interface{})
	variables["empl_load_id"] = emplLoadID
	variables["package_loaded"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("PackageLoaded").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// VehicleStarted send mesage to zeebe engine
func VehicleStarted(orderID, emplDriverID uint) error {
	variables := make(map[string]interface{})
	variables["empl_driver_id"] = emplDriverID
	variables["package_loaded"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("VehicleStarted").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// VehicleArrived send mesage to zeebe engine
func VehicleArrived(orderID uint) error {
	variables := make(map[string]interface{})

	request, err := zbClient.NewPublishMessageCommand().MessageName("VehicleArrived").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// PackageUnloaded send mesage to zeebe engine
func PackageUnloaded(orderID, emplUnloadID uint) error {
	variables := make(map[string]interface{})
	variables["empl_unload_id"] = emplUnloadID
	variables["package_unloaded"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("PackageUnloaded").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// ShipperReceived send mesage to zeebe engine
func ShipperReceived(orderID uint) error {
	variables := make(map[string]interface{})
	variables["shipper_received"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("PackageUnloaded").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// ShipperCalled send mesage to zeebe engine
func ShipperCalled(orderID uint, timeConfirmed int64) error {
	variables := make(map[string]interface{})
	variables["shipper_called"] = true
	variables["time_confirmed"] = timeConfirmed

	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperCalled").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// ShipperShipped send mesage to zeebe engine
func ShipperShipped(orderID uint) error {
	variables := make(map[string]interface{})

	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperShipped").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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

// ShipperConfirmed send mesage to zeebe engine
func ShipperConfirmed(orderID uint, shipperConfirmed string) error {
	variables := make(map[string]interface{})
	variables["shipper_confirmed"] = shipperConfirmed

	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperConfirmed").CorrelationKey(fmt.Sprint(orderID)).VariablesFromMap(variables)
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
