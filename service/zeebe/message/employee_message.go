package message

import (
	"context"
	"fmt"
	"os"

	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

var (
	// ZbClient client to connect with zeebe engine
	zbClient zbc.Client
)

// ConnectZeebeEngine function
func ConnectZeebeEngine() {
	gatewayAddress := os.Getenv("BROKER_ADDRESS")
	newZbClient, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress:         gatewayAddress,
		UsePlaintextConnection: true,
	})

	if err != nil {
		panic(err)
	}

	zbClient = newZbClient
}

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

// MoneyReceived2 send mesage to zeebe engine
func MoneyReceived2(payEmployeeID, orderID uint) error {
	variables := make(map[string]interface{})
	variables["pay_employee_id"] = payEmployeeID
	variables["pay_status"] = true

	ctx := context.Background()
	request, err := zbClient.NewSetVariablesCommand().ElementInstanceKey(2251799813685251).VariablesFromMap(variables)
	_, err = request.Send(ctx)
	if err != nil {
		return err
	}
	_, err = zbClient.NewResolveIncidentCommand().IncidentKey(2251799813685251).Send(ctx)
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
func PackageLoaded(emplLoadID, orderID uint) error {
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
func VehicleStarted(emplDriverID, orderID uint) error {
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
func PackageUnloaded(emplUnloadID, orderID uint) error {
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
func ShipperCalled(timeConfirmed int64, orderID uint) error {
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
func ShipperConfirmed(shipperConfirmed string, orderID uint) error {
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
