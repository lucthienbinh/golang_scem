package workflow

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lucthienbinh/golang_scem/internal/model"
	"github.com/zeebe-io/zeebe/clients/go/pkg/zbc"
)

var (
	// ZbClient client to connect with zeebe engine
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

// DeployFullShipWorkflow function
func DeployFullShipWorkflow() {

	ctx := context.Background()
	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile(os.Getenv("FULL_SHIP_ZB_FILE_1")).Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())
}

// DeployLongShipWorkflow function
func DeployLongShipWorkflow() {

	ctx := context.Background()
	response, err := zbClient.NewDeployWorkflowCommand().AddResourceFile(os.Getenv("LONG_SHIP_ZB_FILE_1")).Send(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(response.String())
}

// CreateFullShipInstance of workflow
func CreateFullShipInstance(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["order_id"] = orderWorkflowData.OrderID
	variables["pay_method"] = orderWorkflowData.PayMethod
	variables["shipper_receive_money"] = orderWorkflowData.ShipperReceiveMoney
	variables["use_long_ship"] = orderWorkflowData.UseLongShip
	variables["use_short_ship"] = orderWorkflowData.UseShortShip
	variables["customer_receive_id"] = orderWorkflowData.CustomerReceiveID

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(os.Getenv("FULL_SHIP_ZB_ID_1")).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		return "", uint(0), err
	}

	ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		return "", uint(0), err
	}
	log.Println(msg.String())
	return strconv.Itoa(int(msg.WorkflowKey)), uint(msg.WorkflowInstanceKey), nil
}

// CreateLongShipInstance of workflow
func CreateLongShipInstance(orderWorkflowData *model.OrderWorkflowData) (string, uint, error) {

	// After the workflow is deployed.
	variables := make(map[string]interface{})
	variables["long_ship_id"] = orderWorkflowData.OrderID

	request, err := zbClient.NewCreateInstanceCommand().BPMNProcessId(os.Getenv("LONG_SHIP_ZB_ID_1")).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		return "", uint(0), err
	}

	ctx := context.Background()
	msg, err := request.Send(ctx)
	if err != nil {
		return "", uint(0), err
	}
	log.Println(msg.String())
	return strconv.Itoa(int(msg.WorkflowKey)), uint(msg.WorkflowInstanceKey), nil
}
