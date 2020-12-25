package message

import (
	"context"
	"fmt"

	ZBWorkflow "github.com/lucthienbinh/golang_scem/service/zeebe/workflow"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

var (
	zbClient zbc.Client
)

// ImportZeebeClientInstance to this package
func ImportZeebeClientInstance() {
	zbClient = ZBWorkflow.ZbClient
}

// MoneyReceived send mesage to zeebe engine
func MoneyReceived(payEmployeeID, orderID uint) error {
	variables := make(map[string]interface{})
	variables["order_id"] = orderID
	variables["pay_employee_id"] = payEmployeeID
	variables["pay_status"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("MoneyReceived").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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

// PaymentInfoUpdated send mesage to zeebe engine
func PaymentInfoUpdated(orderID uint) error {
	variables := make(map[string]interface{})

	request, err := zbClient.NewPublishMessageCommand().MessageName("PaymentInfoUpdated").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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
func PackageLoaded(emplLoadID, orderID uint) error {
	variables := make(map[string]interface{})
	variables["empl_load_id"] = emplLoadID
	variables["package_loaded"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("PackageLoaded").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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
func VehicleStarted(emplDriverID, orderID uint) error {
	variables := make(map[string]interface{})
	variables["empl_driver_id"] = emplDriverID
	variables["package_loaded"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("VehicleStarted").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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

	request, err := zbClient.NewPublishMessageCommand().MessageName("VehicleArrived").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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
func PackageUnloaded(emplUnloadID, orderID uint) error {
	variables := make(map[string]interface{})
	variables["empl_unload_id"] = emplUnloadID
	variables["package_unloaded"] = true

	request, err := zbClient.NewPublishMessageCommand().MessageName("PackageUnloaded").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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

	request, err := zbClient.NewPublishMessageCommand().MessageName("PackageUnloaded").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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
func ShipperCalled(timeConfirmed int64, orderID uint) error {
	variables := make(map[string]interface{})
	variables["shipper_called"] = true
	variables["time_confirmed"] = timeConfirmed

	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperCalled").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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

	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperShipped").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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

// ShipperConfirmed send mesage to zeebe engine
func ShipperConfirmed(shipperConfirmed string, orderID uint) error {
	variables := make(map[string]interface{})
	variables["shipper_confirmed"] = shipperConfirmed

	request, err := zbClient.NewPublishMessageCommand().MessageName("ShipperConfirmed").CorrelationKey(fmt.Sprintf("%x", orderID)).VariablesFromMap(variables)
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
